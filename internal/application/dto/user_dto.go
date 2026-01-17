package dto

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Status int    `json:"status" binding:"omitempty,oneof=1 2"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	TeamID uint64 `json:"team_id"`
	Role   string `json:"role"`
	Status int    `json:"status"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	List []*UserResponse `json:"list"`
	Page PageResponse    `json:"page"`
}
