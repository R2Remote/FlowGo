package dao

type TeamPO struct {
	BasePO
	Name        string `gorm:"column:name;not null"`
	Description string `gorm:"column:description"`
	OwnerId     uint64 `gorm:"column:owner_id"`
}

func (TeamPO) TableName() string {
	return "teams"
}
