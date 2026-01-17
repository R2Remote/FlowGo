package repository

import (
	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/domain/repository"
	"FLOWGO/internal/infrastructure/dao"
	"context"
	"time"

	"gorm.io/gorm"
)

// projectsRepositoryImpl 项目仓储实现
type projectsRepository struct {
	db *gorm.DB
}

// NewProjectsRepository 创建项目仓储实例
func NewProjectsRepository(db *gorm.DB) repository.ProjectsRepository {
	return &projectsRepository{db: db}
}

// Create 创建项目
func (r *projectsRepository) Create(ctx context.Context, project *entity.Project) error {
	po := r.toPO(project)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	// 回写ID和时间
	project.ID = po.ID
	project.CreatedAt = po.CreatedAt
	project.UpdatedAt = po.UpdatedAt
	return nil
}

// Delete 删除项目
func (r *projectsRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&dao.ProjectPO{}, id).Error
}

// FindById 根据ID查找
func (r *projectsRepository) FindByID(ctx context.Context, id uint64) (*entity.Project, error) {
	var po dao.ProjectPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.toEntity(&po), nil
}

// List 列表查询
func (r *projectsRepository) List(ctx context.Context, page, pageSize int) ([]*entity.Project, int64, error) {
	var pos []*dao.ProjectPO
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数
	// GORM with gorm.DeletedAt handles soft delete check automatically
	if err := r.db.WithContext(ctx).Model(&dao.ProjectPO{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域对象
	projects := make([]*entity.Project, len(pos))
	for i, po := range pos {
		projects[i] = r.toEntity(po)
	}

	return projects, total, nil
}

// Update 更新项目
func (r *projectsRepository) Update(ctx context.Context, project *entity.Project) error {
	po := r.toPO(project)
	// BasePO updates
	// Ensure ID is set for update
	return r.db.WithContext(ctx).Model(&dao.ProjectPO{BasePO: dao.BasePO{ID: project.ID}}).Updates(po).Error
}

// ListAvailableTeams 列表查询可用团队
func (r *projectsRepository) ListAvailableTeams(ctx context.Context) ([]*entity.Team, error) {
	var pos []*dao.TeamPO
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	teams := make([]*entity.Team, len(pos))
	for i, po := range pos {
		teams[i] = &entity.Team{
			BaseEntity: entity.BaseEntity{
				ID:        po.ID,
				CreatedAt: po.CreatedAt,
				UpdatedAt: po.UpdatedAt,
				DeletedAt: &po.DeletedAt.Time,
			},
			Name:        po.Name,
			Description: po.Description,
			OwnerId:     po.OwnerId,
		}
	}
	return teams, nil
}

// AddTeams 添加团队
func (r *projectsRepository) AddTeams(ctx context.Context, projectId uint64, teamIds []uint64) error {
	if len(teamIds) == 0 {
		return nil
	}
	pos := make([]*dao.ProjectTeamPO, 0, len(teamIds))
	for _, teamId := range teamIds {
		pos = append(pos, &dao.ProjectTeamPO{
			ProjectId: projectId,
			TeamId:    teamId,
		})
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

// DeleteTeamsByProjectId 删除团队
func (r *projectsRepository) DeleteTeamsByProjectId(ctx context.Context, projectId uint64) error {
	return r.db.WithContext(ctx).Delete(&dao.ProjectTeamPO{}, "project_id = ?", projectId).Error
}

// ListTeamsByProjectId 列表查询团队
func (r *projectsRepository) ListTeamIdsByProjectId(ctx context.Context, projectId uint64) ([]uint64, error) {
	var pos []*dao.ProjectTeamPO
	err := r.db.WithContext(ctx).Where("project_id = ?", projectId).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	teamIds := make([]uint64, len(pos))
	for i, po := range pos {
		teamIds[i] = po.TeamId
	}
	return teamIds, nil
}

// ListUsersByProjectId 根据项目ID列表查询用户
func (r *projectsRepository) ListUsersByProjectId(ctx context.Context, projectId uint64) ([]*entity.User, error) {
	// 只查询 projects_users 表，不再通过 Team 关联
	var userIDs []uint64

	err := r.db.WithContext(ctx).
		Model(&dao.ProjectUserPO{}).
		Where("project_id = ?", projectId).
		Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return []*entity.User{}, nil
	}

	// 3. 去重 (虽然 projects_users 应该是唯一的，但为了保险)
	uniqueIDs := make(map[uint64]bool)
	var distinctIDs []uint64
	for _, id := range userIDs {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			distinctIDs = append(distinctIDs, id)
		}
	}

	// 4. 查询用户详情
	var pos []*dao.UserPO
	err = r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", distinctIDs).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		users[i] = &entity.User{
			BaseEntity: entity.BaseEntity{
				ID:        po.ID,
				CreatedAt: po.CreatedAt,
				UpdatedAt: po.UpdatedAt,
				DeletedAt: &po.DeletedAt.Time,
			},
			Name:   po.Name,
			Email:  po.Email,
			Status: 1, // 数据库中是字符串(online/offline)，暂时映射为1(正常)
			Avatar: po.Avatar,
			TeamID: po.TeamID,
			Role:   po.Role,
		}
	}
	return users, nil
}

// ListAvailableUsers 列表查询可用用户
func (r *projectsRepository) ListAvailableUsers(ctx context.Context) ([]*entity.User, error) {
	var pos []*dao.UserPO
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		users[i] = &entity.User{
			BaseEntity: entity.BaseEntity{
				ID:        po.ID,
				CreatedAt: po.CreatedAt,
				UpdatedAt: po.UpdatedAt,
				DeletedAt: &po.DeletedAt.Time,
			},
			Name:   po.Name,
			Email:  po.Email,
			Status: 1, // 数据库中是字符串(online/offline)，暂时映射为1(正常)
		}
	}
	return users, nil
}

// AddUsers 添加项目成员
func (r *projectsRepository) AddUsers(ctx context.Context, projectId uint64, userIds []uint64) error {
	if len(userIds) == 0 {
		return nil
	}
	// 简单实现：尝试批量插入，利用唯一索引忽略重复或在应用层过滤
	// 这里选择应用层过滤以避免错误
	var existingUserIds []uint64
	r.db.WithContext(ctx).Model(&dao.ProjectUserPO{}).
		Where("project_id = ? AND user_id IN ?", projectId, userIds).
		Pluck("user_id", &existingUserIds)

	existingMap := make(map[uint64]bool)
	for _, id := range existingUserIds {
		existingMap[id] = true
	}

	var newPOs []*dao.ProjectUserPO
	for _, uid := range userIds {
		if !existingMap[uid] {
			newPOs = append(newPOs, &dao.ProjectUserPO{
				ProjectId: projectId,
				UserId:    uid,
				Role:      "member",
			})
		}
	}

	if len(newPOs) > 0 {
		return r.db.WithContext(ctx).Create(&newPOs).Error
	}

	for _, uid := range userIds {
		if !existingMap[uid] {
			newPOs = append(newPOs, &dao.ProjectUserPO{
				ProjectId: projectId,
				UserId:    uid,
				Role:      "member",
			})
		}
	}

	if len(newPOs) > 0 {
		return r.db.WithContext(ctx).Create(&newPOs).Error
	}
	return nil
}

// ListUsersInProjectTeams 获取项目关联团队下的所有用户 (用于添加成员时的候选列表)
func (r *projectsRepository) ListUsersInProjectTeams(ctx context.Context, projectId uint64) ([]*entity.User, error) {
	var teamIDs []uint64
	// 1. 先查项目关联的团队ID
	err := r.db.WithContext(ctx).
		Model(&dao.ProjectTeamPO{}).
		Where("project_id = ?", projectId).
		Pluck("team_id", &teamIDs).Error
	if err != nil {
		return nil, err
	}

	if len(teamIDs) == 0 {
		return []*entity.User{}, nil
	}

	// 2. 再查这些团队下的用户
	var pos []*dao.UserPO
	err = r.db.WithContext(ctx).
		Model(&dao.UserPO{}).
		Where("team_id IN ? AND deleted_at IS NULL", teamIDs).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		users[i] = &entity.User{
			BaseEntity: entity.BaseEntity{
				ID:        po.ID,
				CreatedAt: po.CreatedAt,
				UpdatedAt: po.UpdatedAt,
				DeletedAt: &po.DeletedAt.Time,
			},
			Name:   po.Name,
			Email:  po.Email,
			Status: 1,
			Avatar: po.Avatar,
			TeamID: po.TeamID,
			Role:   po.Role,
		}
	}
	return users, nil
}

// RemoveUsers 移除项目成员
func (r *projectsRepository) RemoveUsers(ctx context.Context, projectId uint64, userId uint64) error {
	return r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectId, userId).
		Delete(&dao.ProjectUserPO{}).Error
}

// Helper methods

func (r *projectsRepository) toPO(e *entity.Project) *dao.ProjectPO {
	var deadline *time.Time
	if !e.Deadline.IsZero() {
		deadline = &e.Deadline
	}
	var startDate *time.Time
	if !e.StartDate.IsZero() {
		startDate = &e.StartDate
	}

	return &dao.ProjectPO{
		BasePO: dao.BasePO{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			// DeletedAt is not manually set usually for updates/creates unless specific logic
		},
		Name:        e.Name,
		Description: e.Description,
		OwnerId:     e.OwnerID,
		Status:      int(e.Status),
		Deadline:    deadline,
		StartDate:   startDate,
		Progress:    e.Progress,
		Priority:    int(e.Priority),
		CoverImage:  e.CoverImage,
	}
}

func (r *projectsRepository) toEntity(po *dao.ProjectPO) *entity.Project {
	var deadline time.Time
	if po.Deadline != nil {
		deadline = *po.Deadline
	}
	var startDate time.Time
	if po.StartDate != nil {
		startDate = *po.StartDate
	}

	e := &entity.Project{
		BaseEntity: entity.BaseEntity{
			ID:        po.ID,
			CreatedAt: po.CreatedAt,
			UpdatedAt: po.UpdatedAt,
		},
		Name:        po.Name,
		Description: po.Description,
		OwnerID:     po.OwnerId,
		Status:      entity.ProjectStatus(po.Status),
		Deadline:    deadline,
		StartDate:   startDate,
		Progress:    po.Progress,
		Priority:    entity.ProjectPriority(po.Priority),
		CoverImage:  po.CoverImage,
	}
	if po.DeletedAt.Valid {
		e.DeletedAt = &po.DeletedAt.Time
	}
	return e
}
