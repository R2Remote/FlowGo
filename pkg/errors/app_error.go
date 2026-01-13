package errors

import "fmt"

// AppError 应用错误
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 常用错误
var (
	ErrNotFound     = NewAppError(404, "资源未找到", nil)
	ErrUnauthorized = NewAppError(401, "未授权", nil)
	ErrForbidden    = NewAppError(403, "禁止访问", nil)
	ErrBadRequest   = NewAppError(400, "请求参数错误", nil)
)
