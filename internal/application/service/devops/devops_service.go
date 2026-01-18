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
	"strings"
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

// ConfigRepo 配置全局仓库
func (s *DevOpsService) ConfigRepo(ctx context.Context, req dto.ConfigRepoRequest) (*dto.ConfigRepoResponse, error) {
	// 1. 检查是否存在配置，如果存在则更新，否则创建
	config, err := s.devopsRepo.GetRepoConfig(ctx)
	if err != nil {
		return nil, errors.New("查询仓库配置失败")
	}

	if config == nil {
		config = &devops.RepoConfig{
			// Global, no ProjectID
		}
		// 生成 Webhook Secret
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err == nil {
			config.WebhookSecret = hex.EncodeToString(bytes)
		}
	}

	config.Type = devops.RepoType(req.Type)
	config.RepoURL = req.RepoURL
	if req.AccessToken != "" {
		config.AccessToken = req.AccessToken // 实际场景应加密存储
	}

	if err := s.devopsRepo.SaveRepoConfig(ctx, config); err != nil {
		return nil, errors.New("保存仓库配置失败")
	}

	return &dto.ConfigRepoResponse{
		ID:      config.ID,
		RepoURL: config.RepoURL,
	}, nil
}

// GetSummary 获取 DevOps 概览
func (s *DevOpsService) GetSummary(ctx context.Context) (*dto.DevOpsSummaryResponse, error) {
	config, err := s.devopsRepo.GetRepoConfig(ctx)
	if err != nil {
		return nil, errors.New("查询仓库配置失败")
	}

	var repoResp *dto.ConfigRepoResponse
	var pipelineResps []*dto.PipelineRecordResponse

	if config != nil {
		repoResp = &dto.ConfigRepoResponse{
			ID:      config.ID,
			RepoURL: config.RepoURL,
		}

		// 查询最近的流水线记录
		records, err := s.devopsRepo.ListPipelineRecords(ctx, 10)
		if err != nil {
			return nil, errors.New("查询流水线记录失败")
		}

		for _, r := range records {
			pipelineResps = append(pipelineResps, &dto.PipelineRecordResponse{
				ID:         r.ID,
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
	}

	return &dto.DevOpsSummaryResponse{
		RepoConfig: repoResp,
		Pipelines:  pipelineResps,
	}, nil
}

// HandleWebhook 处理 Webhook 事件
func (s *DevOpsService) HandleWebhook(ctx context.Context, payload dto.WebhookPayload) error {
	// 根据 RepoURL 找到对应的 RepoConfig
	config, err := s.devopsRepo.GetRepoConfigByRepoURL(ctx, payload.RepoURL)
	if err != nil {
		return errors.New("查询仓库配置失败")
	}

	// 默认 ID 1
	var repoConfigID uint64 = 1
	if config != nil {
		repoConfigID = config.ID
	} else {
		// 如果未找到配置，但可以通过 URL 识别项目，则继续处理（使用默认 ID 记录日志）
		// 这允许用户不配数据库也能跑自动部署
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
		RepoConfigID: repoConfigID,
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

		var deployType string
		// 简单的关键词匹配逻辑 (Pragmatic approach)
		// 实际应该更严谨，或者在 Config 里配置 "DeployScript"
		if strings.Contains(payload.RepoURL, "FlowBoard") {
			deployType = "frontend"
		} else if strings.Contains(payload.RepoURL, "FlowGo") {
			deployType = "backend"
		}

		if deployType != "" {
			go func() {
				// Pass deployType to TriggerDeployment
				if err := s.TriggerDeployment(context.Background(), deployType); err != nil {
					_ = err
				}
			}()
		}
	}

	return s.devopsRepo.SavePipelineRecord(ctx, record)
}

// TriggerDeployment 触发部署
func (s *DevOpsService) TriggerDeployment(ctx context.Context, deployType string) error {
	scriptName := "scripts/deploy_backend.sh" // Default
	commitMsg := "Manual Deployment Trigger"

	if deployType == "frontend" {
		scriptName = "scripts/deploy_frontend.sh"
		commitMsg = "Frontend Deployment"
	} else if deployType == "backend" {
		scriptName = "scripts/deploy_backend.sh"
		commitMsg = "Backend Deployment"
	}

	// 1. 创建一条新的流水线记录
	record := &devops.PipelineRecord{
		RepoConfigID: 1, // Global Config ID
		Status:       devops.PipelineStatusRunning,
		CommitMsg:    commitMsg,
		Author:       "System",
		StartedAt:    nowPtr(),
	}

	// Try to get actual config (Just for logging purpose, reuse global one)
	config, _ := s.devopsRepo.GetRepoConfig(ctx)
	if config != nil {
		record.RepoConfigID = config.ID
	}

	if err := s.devopsRepo.SavePipelineRecord(ctx, record); err != nil {
		return err
	}

	// 2. 异步执行脚本
	go func(recordID uint64, script string) {
		bgCtx := context.Background()

		// Execute shell script
		cmd := exec.Command("bash", script)
		// In production, use absolute path

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
