package dao

import (
	"time"

	"gorm.io/gorm"
)

// BasePO 基础持久化对象
type BasePO struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
