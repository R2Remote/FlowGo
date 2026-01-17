package devops

import (
	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/entity/devops"
	"FLOWGO/internal/domain/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

	// 创建流水线记录
	record := &devops.PipelineRecord{
		// Global, no ProjectID
		RepoConfigID: config.ID,
		ExternalID:   payload.CommitSHA, // 简化：暂时用 SHA 作为 ExternalID，如果是 Pipeline Event 应该用 RunID
		Ref:          payload.Ref,
		CommitSHA:    payload.CommitSHA,
		CommitMsg:    payload.CommitMsg,
		Author:       payload.Author,
		Status:       status,
		// Duration, StartedAt 等字段需根据具体 Payload 填充，此处简化
	}

	return s.devopsRepo.SavePipelineRecord(ctx, record)
}
