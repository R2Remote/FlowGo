package handler

import (
	"github.com/gin-gonic/gin"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/application/service"
	apperrors "FLOWGO/pkg/errors"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	BaseHandler
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，返回JWT Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "登录信息"
// @Success 200 {object} dto.Response{data=dto.LoginResponse}
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			h.HandleError(c, appErr.Code, appErr.Message)
		} else {
			h.HandleInternalError(c, err.Error())
		}
		return
	}

	h.HandleSuccess(c, result)
}
