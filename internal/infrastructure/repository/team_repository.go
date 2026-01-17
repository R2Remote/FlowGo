package repository

import (
	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/domain/repository"
	"FLOWGO/internal/infrastructure/dao"
	"context"

	"gorm.io/gorm"
)

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) repository.TeamRepository {
	return &teamRepository{
		db: db,
	}
}

// ListAvailableTeams 列表查询可用团队
func (r *teamRepository) ListAvailableTeams(ctx context.Context) ([]*entity.Team, error) {
	var pos []*dao.TeamPO
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	teams := make([]*entity.Team, len(pos))
	for i, po := range pos {
		teams[i] = &entity.Team{
			BaseEntity: entity.BaseEntity{
				ID:        po.ID,
				CreatedAt: po.CreatedAt,
				UpdatedAt: po.UpdatedAt,
				DeletedAt: &po.DeletedAt.Time, // Assuming valid given query
			},
			Name:        po.Name,
			Description: po.Description,
			OwnerId:     po.OwnerId,
		}
	}
	return teams, nil
}
