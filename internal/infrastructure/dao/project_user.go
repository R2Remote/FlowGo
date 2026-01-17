package dao

// ProjectUserPO 项目用户关联表
type ProjectUserPO struct {
	BasePO
	ProjectId uint64 `json:"project_id"`
	UserId    uint64 `json:"user_id"`
	Role      string `json:"role"`
}

func (ProjectUserPO) TableName() string {
	return "projects_users"
}
