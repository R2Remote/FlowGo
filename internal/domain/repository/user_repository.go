package repository

import (
	"context"
	"FLOWGO/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[entity.User]
	
	// FindByUsername 根据用户名查找
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// FindByEmail 根据邮箱查找
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
