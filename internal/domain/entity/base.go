package entity

import "time"

// BaseEntity 基础实体，包含通用字段
type BaseEntity struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// IsDeleted 检查实体是否已删除
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt != nil
}
