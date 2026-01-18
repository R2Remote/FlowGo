package devops

import (
	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/entity/devops"
	"FLOWGO/internal/domain/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os/exec"
	"time"
)

type DevOpsService struct {
	devopsRepo repository.DevOpsRepository
}

func NewDevOpsService(devopsRepo repository.DevOpsRepository) *DevOpsService {
	return &DevOpsService{
		devopsRepo: devopsRepo,
	}
}

// ConfigRepo 配置仓库
func (s *DevOpsService) ConfigRepo(ctx context.Context, req dto.ConfigRepoRequest) (*dto.ConfigRepoResponse, error) {
	var config *devops.RepoConfig
	var err error

	if req.ID > 0 {
		config, err = s.devopsRepo.GetRepoConfig(ctx, req.ID)
		if err != nil {
			return nil, errors.New("查询仓库配置失败")
		}
		if config == nil {
			return nil, errors.New("仓库配置不存在")
		}
	} else {
		config = &devops.RepoConfig{
			// New record
		}
		// 生成 Webhook Secret
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err == nil {
			config.WebhookSecret = hex.EncodeToString(bytes)
		}
	}

	config.Name = req.Name
	config.Type = devops.RepoType(req.Type)
	config.RepoURL = req.RepoURL
	config.DeployScript = req.DeployScript
	if req.AccessToken != "" {
		config.AccessToken = req.AccessToken // 实际场景应加密存储
	}

	if err := s.devopsRepo.SaveRepoConfig(ctx, config); err != nil {
		return nil, errors.New("保存仓库配置失败")
	}

	return &dto.ConfigRepoResponse{
		ID:           config.ID,
		Name:         config.Name,
		RepoURL:      config.RepoURL,
		DeployScript: config.DeployScript,
	}, nil
}

// GetSummary 获取 DevOps 概览
func (s *DevOpsService) GetSummary(ctx context.Context) (*dto.DevOpsSummaryResponse, error) {
	configs, err := s.devopsRepo.ListRepoConfigs(ctx)
	if err != nil {
		return nil, errors.New("查询仓库配置失败")
	}

	var serviceResps []*dto.ConfigRepoResponse
	configMap := make(map[uint64]string)

	for _, c := range configs {
		serviceResps = append(serviceResps, &dto.ConfigRepoResponse{
			ID:           c.ID,
			Name:         c.Name,
			RepoURL:      c.RepoURL,
			DeployScript: c.DeployScript,
		})
		configMap[c.ID] = c.Name
	}

	var pipelineResps []*dto.PipelineRecordResponse
	// 查询最近的流水线记录
	records, err := s.devopsRepo.ListPipelineRecords(ctx, 10)
	if err != nil {
		return nil, errors.New("查询流水线记录失败")
	}

	for _, r := range records {
		repoName := configMap[r.RepoConfigID]
		if repoName == "" {
			repoName = "Unknown"
		}
		pipelineResps = append(pipelineResps, &dto.PipelineRecordResponse{
			ID:         r.ID,
			RepoName:   repoName,
			Status:     string(r.Status),
			Ref:        r.Ref,
			CommitSHA:  r.CommitSHA,
			CommitMsg:  r.CommitMsg,
			Author:     r.Author,
			Duration:   r.Duration,
			StartedAt:  r.StartedAt,
			FinishedAt: r.FinishedAt,
			CreatedAt:  r.CreatedAt,
		})
	}

	return &dto.DevOpsSummaryResponse{
		Services:  serviceResps,
		Pipelines: pipelineResps,
	}, nil
}

// HandleWebhook 处理 Webhook 事件
func (s *DevOpsService) HandleWebhook(ctx context.Context, payload dto.WebhookPayload) error {
	// 根据 RepoURL 找到对应的 RepoConfig
	config, err := s.devopsRepo.GetRepoConfigByRepoURL(ctx, payload.RepoURL)
	if err != nil {
		return errors.New("查询仓库配置失败")
	}
	if config == nil {
		// 未配置该仓库，忽略事件
		return nil
	}

	// 转换状态
	var status devops.PipelineStatus
	switch payload.Status {
	case "success":
		status = devops.PipelineStatusSuccess
	case "failed", "failure":
		status = devops.PipelineStatusFailed
	case "running", "pending":
		status = devops.PipelineStatusRunning
	case "canceled":
		status = devops.PipelineStatusCanceled
	default:
		status = devops.PipelineStatusPending
	}

	// 创建流水线记录 (Logging the event)
	record := &devops.PipelineRecord{
		RepoConfigID: config.ID,
		ExternalID:   payload.CommitSHA,
		Ref:          payload.Ref,
		CommitSHA:    payload.CommitSHA,
		CommitMsg:    payload.CommitMsg,
		Author:       payload.Author,
		Status:       status,
	}

	// 自动触发部署 (仅当 Push 到 main 分支且状态为 success 时)
	if status == devops.PipelineStatusSuccess &&
		(payload.Ref == "refs/heads/main" || payload.Ref == "main") {

		// 直接使用 Config 里的 DeployScript
		if config.DeployScript != "" {
			go func() {
				// Pass ID to TriggerDeployment
				if err := s.TriggerDeployment(context.Background(), config.ID); err != nil {
					_ = err
				}
			}()
		}
	}

	return s.devopsRepo.SavePipelineRecord(ctx, record)
}

// TriggerDeployment 触发部署
func (s *DevOpsService) TriggerDeployment(ctx context.Context, repoConfigID uint64) error {
	// Fetch config to get script path
	config, err := s.devopsRepo.GetRepoConfig(ctx, repoConfigID)
	if err != nil || config == nil {
		return errors.New("配置不存在")
	}

	if config.DeployScript == "" {
		return errors.New("未配置部署脚本")
	}

	scriptName := config.DeployScript
	commitMsg := "Manual Deployment Trigger"

	// 1. 创建一条新的流水线记录
	record := &devops.PipelineRecord{
		RepoConfigID: config.ID,
		Status:       devops.PipelineStatusRunning,
		CommitMsg:    commitMsg,
		Author:       "System",
		StartedAt:    nowPtr(),
	}

	if err := s.devopsRepo.SavePipelineRecord(ctx, record); err != nil {
		return err
	}

	// 2. 异步执行脚本
	go func(recordID uint64, script string) {
		bgCtx := context.Background()

		// Execute shell script
		// Use bash explicitely
		cmd := exec.Command("bash", script)

		output, err := cmd.CombinedOutput()

		endTime := time.Now()
		duration := int64(time.Since(*record.StartedAt).Seconds())

		status := devops.PipelineStatusSuccess
		if err != nil {
			status = devops.PipelineStatusFailed
		}

		updateObj := &devops.PipelineRecord{
			ID:         recordID,
			Status:     status,
			Duration:   duration,
			FinishedAt: &endTime,
			CommitMsg:  "Output: " + string(output),
		}
		if len(updateObj.CommitMsg) > 500 {
			updateObj.CommitMsg = updateObj.CommitMsg[:500] + "..."
		}

		_ = s.devopsRepo.SavePipelineRecord(bgCtx, updateObj)

	}(record.ID, scriptName)

	return nil
}

func nowPtr() *time.Time {
	t := time.Now()
	return &t
}
