package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/application/service"
	apperrors "FLOWGO/pkg/errors"
)

// UserHandler 用户处理器
type UserHandler struct {
	BaseHandler
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "用户信息"
// @Success 200 {object} dto.Response{data=dto.UserResponse}
// @Failure 400 {object} dto.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, user)
}

// GetUser 获取用户
// @Summary 获取用户
// @Description 根据ID获取用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response{data=dto.UserResponse}
// @Failure 404 {object} dto.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.HandleBadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, user)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.Response{data=dto.UserListResponse}
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	result, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccessWithPage(c, result.List, result.Page.Page, result.Page.PageSize, result.Page.Total)
}
