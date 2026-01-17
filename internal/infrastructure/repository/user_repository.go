package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"FLOWGO/internal/domain/entity"
	domainRepo "FLOWGO/internal/domain/repository"
)

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &userRepository{
		db: db,
	}
}

// FindByID 根据ID查找
func (r *userRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

// Create 创建
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除（软删除）
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

// List 列表查询
func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
