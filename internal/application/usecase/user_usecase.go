package usecase

import (
	"context"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/domain/repository"
	apperrors "FLOWGO/pkg/errors"
	"FLOWGO/pkg/utils"
)

// CreateUserUseCase 创建用户用例
type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

// NewCreateUserUseCase 创建用例实例
func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute 执行创建用户
func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 检查用户名是否存在
	exists, err := uc.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.NewAppError(500, "检查用户名失败", err)
	}
	if exists {
		return nil, apperrors.NewAppError(400, "用户名已存在", nil)
	}

	// 检查邮箱是否存在
	exists, err = uc.userRepo.ExistsByEmail(ctx, req.Email)
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
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   1,
	}

	// 保存用户
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewAppError(500, "创建用户失败", err)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

// GetUserUseCase 获取用户用例
type GetUserUseCase struct {
	userRepo repository.UserRepository
}

// NewGetUserUseCase 创建用例实例
func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute 执行获取用户
func (uc *GetUserUseCase) Execute(ctx context.Context, id uint64) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(500, "查询用户失败", err)
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

// ListUsersUseCase 用户列表用例
type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

// NewListUsersUseCase 创建用例实例
func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute 执行获取用户列表
func (uc *ListUsersUseCase) Execute(ctx context.Context, req dto.PageRequest) (*dto.UserListResponse, error) {
	users, total, err := uc.userRepo.List(ctx, req.Page, req.GetPageSize())
	if err != nil {
		return nil, apperrors.NewAppError(500, "查询用户列表失败", err)
	}

	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
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
