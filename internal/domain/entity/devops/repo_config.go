package devops

import "time"

// RepoType 定义仓库类型
type RepoType string

const (
	RepoTypeGitHub RepoType = "github"
	RepoTypeGitLab RepoType = "gitlab"
)

// RepoConfig 仓库配置实体
type RepoConfig struct {
	ID            uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null"`
	Type          RepoType  `json:"type" gorm:"type:varchar(20);not null"`
	RepoURL       string    `json:"repo_url" gorm:"type:varchar(255);not null"`
	DeployScript  string    `json:"deploy_script" gorm:"type:varchar(255)"`
	AccessToken   string    `json:"-" gorm:"type:varchar(255)"` // 不在 JSON 中返回，并在存储时加密
	WebhookSecret string    `json:"-" gorm:"type:varchar(100)"` // Webhook 验证密钥
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
