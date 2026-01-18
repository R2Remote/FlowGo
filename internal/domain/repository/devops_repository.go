package repository

import (
	"FLOWGO/internal/domain/entity/devops"
	"context"
)

// DevOpsRepository DevOps 仓储接口
type DevOpsRepository interface {
	// RepoConfig
	// RepoConfig
	SaveRepoConfig(ctx context.Context, config *devops.RepoConfig) error
	GetRepoConfig(ctx context.Context, id uint64) (*devops.RepoConfig, error)
	ListRepoConfigs(ctx context.Context) ([]*devops.RepoConfig, error)
	GetRepoConfigByRepoURL(ctx context.Context, repoURL string) (*devops.RepoConfig, error)
	DeleteRepoConfig(ctx context.Context, id uint64) error

	// PipelineRecord
	SavePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error
	ListPipelineRecords(ctx context.Context, limit int) ([]*devops.PipelineRecord, error)
	GetPipelineRecordByExternalID(ctx context.Context, externalID string) (*devops.PipelineRecord, error)
}
