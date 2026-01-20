package entity

import (
	"time"

	"gorm.io/gorm"
)

// VisitStat 访问统计实体
type VisitStat struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	IP        string         `gorm:"uniqueIndex;size:50;not null" json:"ip"` // IP地址，唯一索引
	Count     int64          `gorm:"default:1" json:"count"`                 // 访问次数
	LastSeen  time.Time      `json:"last_seen"`                              // 最后访问时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
