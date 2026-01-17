package dao

type ProjectTeamPO struct {
	BasePO
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	ProjectId uint64 `gorm:"column:project_id;index;not null"`
	TeamId    uint64 `gorm:"column:team_id;index;not null"`
}

func (ProjectTeamPO) TableName() string {
	return "projects_teams"
}
