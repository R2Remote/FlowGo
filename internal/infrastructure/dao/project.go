package dao

import (
	"time"
)

// ProjectPO 项目持久化对象
type ProjectPO struct {
	BasePO
	Name        string     `gorm:"not null;type:varchar(100)"`
	Description string     `gorm:"not null;type:text"`
	OwnerId     uint64     `gorm:"not null;index"`
	Status      int        `gorm:"type:tinyint;default:1"`
	Deadline    *time.Time `gorm:"type:datetime"`
	StartDate   *time.Time `gorm:"type:datetime"`
	Progress    int        `gorm:"default:0"`
	Priority    int        `gorm:"type:tinyint;default:2"`
	CoverImage  string     `gorm:"type:varchar(255)"`
}

func (ProjectPO) TableName() string {
	return "projects"
}
