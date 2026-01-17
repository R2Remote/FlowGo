package entity

import "time"

// BaseEntity 基础实体，包含通用字段
type BaseEntity struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// IsDeleted 检查实体是否已删除
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt != nil
}
