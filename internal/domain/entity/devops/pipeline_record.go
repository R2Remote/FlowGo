package devops

import "time"

// PipelineStatus 流水线状态
type PipelineStatus string

const (
	PipelineStatusPending  PipelineStatus = "pending"
	PipelineStatusRunning  PipelineStatus = "running"
	PipelineStatusSuccess  PipelineStatus = "success"
	PipelineStatusFailed   PipelineStatus = "failed"
	PipelineStatusCanceled PipelineStatus = "canceled"
)

// PipelineRecord 流水线记录实体
type PipelineRecord struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	RepoConfigID uint64         `json:"repo_config_id" gorm:"index;not null"`
	ExternalID   string         `json:"external_id" gorm:"type:varchar(100);index"` // 外部系统(如GitHub)的 Run ID
	Ref          string         `json:"ref" gorm:"type:varchar(100)"`               // 分支或标签
	CommitSHA    string         `json:"commit_sha" gorm:"type:varchar(40)"`
	CommitMsg    string         `json:"commit_msg" gorm:"type:text"`
	Author       string         `json:"author" gorm:"type:varchar(100)"`
	Status       PipelineStatus `json:"status" gorm:"type:varchar(20);not null"`
	Duration     int64          `json:"duration"` // 秒
	StartedAt    *time.Time     `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
