package repository

import (
	"FLOWGO/internal/domain/entity"
	"context"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[entity.User]

	// FindByEmail 根据邮箱查找
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
