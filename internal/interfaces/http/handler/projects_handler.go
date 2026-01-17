package handler

import (
	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/application/service"
	"FLOWGO/pkg/contextutil"
	apperrors "FLOWGO/pkg/errors"

	"github.com/gin-gonic/gin"
)

type ProjectsHandler struct {
	BaseHandler
	projectService *service.ProjectService
}

func NewProjectsHandler(
	projectService *service.ProjectService,
) *ProjectsHandler {
	return &ProjectsHandler{
		projectService: projectService,
	}
}

func (h *ProjectsHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}
	userID, err := contextutil.GetUserID(c)
	if err != nil {
		h.HandleError(c, 401, "未授权")
		return
	}
	req.OwnerID = uint64(userID)
	project, err := h.projectService.CreateProject(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, project)
}

func (h *ProjectsHandler) UpdateProject(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, err := contextutil.GetUserID(c)
	if err != nil {
		h.HandleError(c, 401, "未授权")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}
	req.OwnerID = userID
	// 将用户ID设置到请求中（如果需要验证权限）
	// 或者直接传递给 UseCase
	project, err := h.projectService.UpdateProject(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, project)
}

func (h *ProjectsHandler) DeleteProject(c *gin.Context) {
	var req dto.DeleteProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	project, err := h.projectService.DeleteProject(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, project)
}

func (h *ProjectsHandler) GetProject(c *gin.Context) {
	var req dto.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	project, err := h.projectService.GetProject(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, project)
}

func (h *ProjectsHandler) ListProjects(c *gin.Context) {
	var req dto.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	projects, err := h.projectService.ListProjects(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccessWithPage(c, projects.List, projects.Page.Page, projects.Page.PageSize, projects.Page.Total)
}

func (h *ProjectsHandler) ProjectTeams(c *gin.Context) {
	teams, err := h.projectService.ProjectTeams(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, teams)
}

func (h *ProjectsHandler) ProjectAvailableUsers(c *gin.Context) {
	var req dto.ProjectUsersRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}
	users, err := h.projectService.GetProjectAvailableUsers(c.Request.Context(), req.ID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, users)
}

func (h *ProjectsHandler) AddProjectUsers(c *gin.Context) {
	var uriReq struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uriReq); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	var req dto.AddProjectUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	resp, err := h.projectService.AddProjectUsers(c.Request.Context(), req, uriReq.ID)
	if err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}

	h.HandleSuccess(c, resp)
}

func (h *ProjectsHandler) RemoveProjectUser(c *gin.Context) {
	var uriReq struct {
		ID     uint64 `uri:"id" binding:"required"`
		UserID uint64 `uri:"uid" binding:"required"`
	}
	if err := c.ShouldBindUri(&uriReq); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	err := h.projectService.RemoveProjectUser(c.Request.Context(), uriReq.ID, uriReq.UserID)
	if err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}

	h.HandleSuccess(c, nil)
}
