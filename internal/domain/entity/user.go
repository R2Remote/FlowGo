package entity

// User 用户实体（示例）
type User struct {
	BaseEntity
	Name     string `json:"name" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Status   int    `json:"status" gorm:"default:1"` // 1:正常 2:禁用
	Avatar   string `json:"avatar"`
	TeamID   uint64 `json:"team_id"`
	Role     string `json:"role"`
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == 1 && !u.IsDeleted()
}
