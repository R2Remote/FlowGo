package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"FLOWGO/internal/application/dto"
)

// BaseHandler 基础处理器
type BaseHandler struct{}

// HandleSuccess 处理成功响应
func (h *BaseHandler) HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, dto.Success(data))
}

// HandleSuccessWithPage 处理带分页的成功响应
func (h *BaseHandler) HandleSuccessWithPage(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	c.JSON(http.StatusOK, dto.SuccessWithPage(data, page, pageSize, total))
}

// HandleError 处理错误响应
func (h *BaseHandler) HandleError(c *gin.Context, code int, message string) {
	httpStatus := http.StatusOK
	if code >= 400 && code < 600 {
		httpStatus = code
	}
	c.JSON(httpStatus, dto.Error(code, message))
}

// HandleBadRequest 处理400错误
func (h *BaseHandler) HandleBadRequest(c *gin.Context, message string) {
	h.HandleError(c, http.StatusBadRequest, message)
}

// HandleInternalError 处理500错误
func (h *BaseHandler) HandleInternalError(c *gin.Context, message string) {
	h.HandleError(c, http.StatusInternalServerError, message)
}
