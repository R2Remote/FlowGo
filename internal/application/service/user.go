package service

import (
	"context"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/domain/repository"
	apperrors "FLOWGO/pkg/errors"
	"FLOWGO/pkg/utils"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
}

// 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (uc *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 检查邮箱是否存在
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewAppError(500, "检查邮箱失败", err)
	}
	if exists {
		return nil, apperrors.NewAppError(400, "邮箱已存在", nil)
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewAppError(500, "密码加密失败", err)
	}

	// 创建用户实体
	user := &entity.User{
		BaseEntity: entity.BaseEntity{
			ID: utils.GenerateID(),
		},
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   1,
	}

	// 保存用户
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewAppError(500, "创建用户失败", err)
	}

	return &dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}, nil
}

// GetUser 获取用户
func (uc *UserService) GetUser(ctx context.Context, id uint64) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(500, "查询用户失败", err)
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	return &dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}, nil
}

// ListUsers 获取用户列表
func (uc *UserService) ListUsers(ctx context.Context, req dto.PageRequest) (*dto.UserListResponse, error) {
	users, total, err := uc.userRepo.List(ctx, req.Page, req.GetPageSize())
	if err != nil {
		return nil, apperrors.NewAppError(500, "查询用户列表失败", err)
	}

	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status,
		})
	}

	return &dto.UserListResponse{
		List: userResponses,
		Page: dto.PageResponse{
			Page:     req.Page,
			PageSize: req.GetPageSize(),
			Total:    total,
		},
	}, nil
}
