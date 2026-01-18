package dto

import "time"

// ConfigRepoRequest 配置仓库请求
type ConfigRepoRequest struct {
	ID           uint64 `json:"id"` // Optional: for update
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=github gitlab"`
	RepoURL      string `json:"repo_url" binding:"required,url"`
	DeployScript string `json:"deploy_script"`
	AccessToken  string `json:"access_token"`
}

// ConfigRepoResponse 配置仓库响应
type ConfigRepoResponse struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	RepoURL      string `json:"repo_url"`
	DeployScript string `json:"deploy_script"`
	WebhookID    string `json:"webhook_id,omitempty"` // 可选：如果自动创建了 webhook
}

// PipelineRecordResponse 流水线记录响应
type PipelineRecordResponse struct {
	ID            uint64     `json:"id"`
	RepoName      string     `json:"repo_name"` // Add repo name for UI
	Status        string     `json:"status"`
	TriggerSource string     `json:"trigger_source"`
	Ref           string     `json:"ref"`
	CommitSHA     string     `json:"commit_sha"`
	CommitMsg     string     `json:"commit_msg"`
	Author        string     `json:"author"`
	Duration      int64      `json:"duration"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

// DevOpsSummaryResponse DevOps 概览响应
type DevOpsSummaryResponse struct {
	Services  []*ConfigRepoResponse     `json:"services"`
	Pipelines []*PipelineRecordResponse `json:"pipelines"`
}

// WebhookPayload 通用 Webhook 负载 (根据不同平台解析后转换为此结构)
type WebhookPayload struct {
	Event     string
	Ref       string
	CommitSHA string
	CommitMsg string
	Author    string
	Status    string
	RepoURL   string
}
