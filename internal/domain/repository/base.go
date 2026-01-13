package repository

import (
	"context"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	// FindByID 根据ID查找
	FindByID(ctx context.Context, id uint64) (*T, error)

	// Create 创建
	Create(ctx context.Context, entity *T) error

	// Update 更新
	Update(ctx context.Context, entity *T) error

	// Delete 删除（软删除）
	Delete(ctx context.Context, id uint64) error

	// List 列表查询
	List(ctx context.Context, page, pageSize int) ([]*T, int64, error)
}
