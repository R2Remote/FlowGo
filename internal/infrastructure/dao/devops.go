package dao

import (
	"time"
)

// RepoConfigPO 仓库配置持久化对象
type RepoConfigPO struct {
	BasePO
	Type          string `gorm:"type:varchar(20);not null"`
	RepoURL       string `gorm:"type:varchar(255);not null"`
	AccessToken   string `gorm:"type:varchar(255)"`
	WebhookSecret string `gorm:"type:varchar(100)"`
}

func (RepoConfigPO) TableName() string {
	return "devops_repo_configs"
}

// PipelineRecordPO 流水线记录持久化对象
type PipelineRecordPO struct {
	BasePO
	RepoConfigID uint64     `gorm:"index;not null"`
	ExternalID   string     `gorm:"type:varchar(100);index"`
	Ref          string     `gorm:"type:varchar(100)"`
	CommitSHA    string     `gorm:"type:varchar(40)"`
	CommitMsg    string     `gorm:"type:text"`
	Author       string     `gorm:"type:varchar(100)"`
	Status       string     `gorm:"type:varchar(20);not null"`
	Duration     int64      `gorm:"default:0"`
	StartedAt    *time.Time `gorm:"default:null"`
	FinishedAt   *time.Time `gorm:"default:null"`
}

func (PipelineRecordPO) TableName() string {
	return "devops_pipeline_records"
}
