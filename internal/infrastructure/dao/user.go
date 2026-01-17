package dao

type UserPO struct {
	BasePO
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	TeamID uint64 `json:"team_id"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

func (UserPO) TableName() string {
	return "users"
}
