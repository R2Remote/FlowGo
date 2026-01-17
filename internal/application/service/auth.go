package service

import (
	"context"
	"time"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/repository"
	apperrors "FLOWGO/pkg/errors"
	"FLOWGO/pkg/jwt"
)

// LoginUseCase 登录用例
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login 登录
func (uc *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// // 根据用户名查找用户
	// user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	// if err != nil {
	// 	return nil, apperrors.NewAppError(500, "查询用户失败", err)
	// }
	// if user == nil {
	// 	return nil, apperrors.NewAppError(401, "用户名或密码错误", nil)
	// }

	// // 检查用户状态
	// if !user.IsActive() {
	// 	return nil, apperrors.NewAppError(403, "用户已被禁用", nil)
	// }

	// // 验证密码
	// if !utils.CheckPassword(req.Password, user.Password) {
	// 	return nil, apperrors.NewAppError(401, "用户名或密码错误", nil)
	// }

	// 生成JWT Token
	token, err := jwt.GenerateToken(1, "test")
	if err != nil {
		return nil, apperrors.NewAppError(500, "生成Token失败", err)
	}

	// 返回登录响应
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:     1,
			Name:   "test",
			Email:  "test@test.com",
			Status: 1,
		},
		ExpiresIn: int64(24 * time.Hour.Seconds()), // 24小时
	}, nil
}
