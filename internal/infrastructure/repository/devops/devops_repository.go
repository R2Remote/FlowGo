package devops

import (
	"FLOWGO/internal/domain/entity/devops"
	"FLOWGO/internal/domain/repository"
	"FLOWGO/internal/infrastructure/dao"
	"context"

	"gorm.io/gorm"
)

type devopsRepository struct {
	db *gorm.DB
}

func NewDevOpsRepository(db *gorm.DB) repository.DevOpsRepository {
	return &devopsRepository{db: db}
}

// SaveRepoConfig 保存仓库配置
func (r *devopsRepository) SaveRepoConfig(ctx context.Context, config *devops.RepoConfig) error {
	po := r.toRepoConfigPO(config)
	if config.ID == 0 {
		return r.db.WithContext(ctx).Create(po).Error
	}
	return r.db.WithContext(ctx).Save(po).Error
}

// GetRepoConfig 获取全局仓库配置 (Singleton)
func (r *devopsRepository) GetRepoConfig(ctx context.Context) (*devops.RepoConfig, error) {
	var po dao.RepoConfigPO
	// Assuming only one record or getting the latest/first
	err := r.db.WithContext(ctx).First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.toRepoConfigEntity(&po), nil
}

// GetRepoConfigByRepoURL 根据仓库URL获取配置
func (r *devopsRepository) GetRepoConfigByRepoURL(ctx context.Context, repoURL string) (*devops.RepoConfig, error) {
	var po dao.RepoConfigPO
	err := r.db.WithContext(ctx).Where("repo_url = ?", repoURL).First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.toRepoConfigEntity(&po), nil
}

// DeleteRepoConfig 删除仓库配置
func (r *devopsRepository) DeleteRepoConfig(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&dao.RepoConfigPO{}, id).Error
}

// SavePipelineRecord 保存流水线记录
func (r *devopsRepository) SavePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error {
	po := r.toPipelineRecordPO(record)

	// 如果 ExternalID 存在，先尝试查找
	if record.ExternalID != "" && record.ID == 0 {
		var existPO dao.PipelineRecordPO
		err := r.db.WithContext(ctx).Where("external_id = ?", record.ExternalID).First(&existPO).Error
		if err == nil {
			po.ID = existPO.ID
			po.CreatedAt = existPO.CreatedAt
		}
	}

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
			return err
		}
		// 回写 ID
		record.ID = po.ID
		record.CreatedAt = po.CreatedAt
		return nil
	}

	return r.db.WithContext(ctx).Model(&dao.PipelineRecordPO{BasePO: dao.BasePO{ID: po.ID}}).Updates(po).Error
}

// ListPipelineRecords 获取流水线记录列表
func (r *devopsRepository) ListPipelineRecords(ctx context.Context, limit int) ([]*devops.PipelineRecord, error) {
	var pos []*dao.PipelineRecordPO
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	records := make([]*devops.PipelineRecord, len(pos))
	for i, po := range pos {
		records[i] = r.toPipelineRecordEntity(po)
	}
	return records, nil
}

// GetPipelineRecordByExternalID 根据外部ID获取记录
func (r *devopsRepository) GetPipelineRecordByExternalID(ctx context.Context, externalID string) (*devops.PipelineRecord, error) {
	var po dao.PipelineRecordPO
	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toPipelineRecordEntity(&po), nil
}

// Converters

func (r *devopsRepository) toRepoConfigPO(e *devops.RepoConfig) *dao.RepoConfigPO {
	return &dao.RepoConfigPO{
		BasePO: dao.BasePO{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		// ProjectID removed
		Type:          string(e.Type),
		RepoURL:       e.RepoURL,
		AccessToken:   e.AccessToken,
		WebhookSecret: e.WebhookSecret,
	}
}

func (r *devopsRepository) toRepoConfigEntity(po *dao.RepoConfigPO) *devops.RepoConfig {
	return &devops.RepoConfig{
		ID: po.ID,
		// ProjectID removed
		Type:          devops.RepoType(po.Type),
		RepoURL:       po.RepoURL,
		AccessToken:   po.AccessToken,
		WebhookSecret: po.WebhookSecret,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

func (r *devopsRepository) toPipelineRecordPO(e *devops.PipelineRecord) *dao.PipelineRecordPO {
	return &dao.PipelineRecordPO{
		BasePO: dao.BasePO{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		// ProjectID removed
		RepoConfigID: e.RepoConfigID,
		ExternalID:   e.ExternalID,
		Ref:          e.Ref,
		CommitSHA:    e.CommitSHA,
		CommitMsg:    e.CommitMsg,
		Author:       e.Author,
		Status:       string(e.Status),
		Duration:     e.Duration,
		StartedAt:    e.StartedAt,
		FinishedAt:   e.FinishedAt,
	}
}

func (r *devopsRepository) toPipelineRecordEntity(po *dao.PipelineRecordPO) *devops.PipelineRecord {
	return &devops.PipelineRecord{
		ID: po.ID,
		// ProjectID removed
		RepoConfigID: po.RepoConfigID,
		ExternalID:   po.ExternalID,
		Ref:          po.Ref,
		CommitSHA:    po.CommitSHA,
		CommitMsg:    po.CommitMsg,
		Author:       po.Author,
		Status:       devops.PipelineStatus(po.Status),
		Duration:     po.Duration,
		StartedAt:    po.StartedAt,
		FinishedAt:   po.FinishedAt,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}
