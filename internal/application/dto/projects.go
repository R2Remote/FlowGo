package dto

import (
	"FLOWGO/pkg/utils"
)

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	OwnerID     uint64 `json:"owner_id"`
}

// CreateProjectResponse 创建项目响应
type CreateProjectResponse struct {
	ID uint64 `json:"id"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ID          uint64     `json:"id" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description" binding:"required"`
	OwnerID     uint64     `json:"owner_id" binding:"required"`
	Status      int        `json:"status" binding:"required"`
	Deadline    utils.Time `json:"deadline" binding:"required"`
	StartDate   utils.Time `json:"start_date" binding:"required"`
	Priority    int        `json:"priority" binding:"required"`
	CoverImage  string     `json:"cover_image" binding:"omitempty"`
	TeamIds     []uint64   `json:"team_ids" binding:"omitempty"`
	Tags        []string   `json:"tags" binding:"omitempty"`
}

// UpdateProjectResponse 更新项目响应
type UpdateProjectResponse struct {
	ID          uint64     `json:"id"`
	TeamIds     []uint64   `json:"team_ids"`
	Tags        []string   `json:"tags"`
	Priority    int        `json:"priority"`
	CoverImage  string     `json:"cover_image"`
	Deadline    utils.Time `json:"deadline"`
	StartDate   utils.Time `json:"start_date"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OwnerID     uint64     `json:"owner_id"`
	Status      int        `json:"status"`
}

// DeleteProjectRequest 删除项目请求
type DeleteProjectRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

// DeleteProjectResponse 删除项目响应
type DeleteProjectResponse struct {
	ID uint64 `json:"id"`
}

// GetProjectRequest 获取项目请求
type GetProjectRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

// GetProjectResponse 获取项目响应
type GetProjectResponse struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	OwnerId     uint64          `json:"owner_id"`
	Status      int             `json:"status"`
	Deadline    utils.Time      `json:"deadline"`
	StartDate   utils.Time      `json:"start_date"`
	Progress    int             `json:"progress"`
	Priority    int             `json:"priority"`
	CoverImage  string          `json:"cover_image"`
	Tags        []string        `json:"tags"`
	TeamIds     []uint64        `json:"team_ids"`
	Users       []*UserResponse `json:"users"`
	CreatedAt   utils.Time      `json:"created_at"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	List []*ProjectResponse `json:"list"`
	Page PageResponse       `json:"page"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OwnerId     uint64     `json:"owner_id"`
	Status      int        `json:"status"`
	Deadline    utils.Time `json:"deadline"`
	StartDate   utils.Time `json:"start_date"`
	Progress    int        `json:"progress"`
	Priority    int        `json:"priority"`
}

// ProjectTeamsResponse 项目团队响应
type ProjectTeamsResponse struct {
	Teams []*TeamResponse `json:"teams"`
}

// TeamResponse 团队响应
type TeamResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ProjectUsersRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

type AddProjectUsersRequest struct {
	Users []uint64 `json:"users" binding:"required"`
}

type ProjectUsersResponse struct {
	Users []*UserResponse `json:"users"`
}
