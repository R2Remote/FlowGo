package service

import (
	"context"
	"errors"

	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/domain/repository"
	"FLOWGO/pkg/utils"
)

type ProjectService struct {
	projectRepo repository.ProjectsRepository
}

func NewProjectService(projectRepo repository.ProjectsRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*dto.CreateProjectResponse, error) {
	project := entity.NewProject(req.Name, req.Description, req.OwnerID)
	err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, errors.New("创建项目失败")
	}
	return &dto.CreateProjectResponse{
		ID: project.ID,
	}, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, req dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error) {
	// 先查找现有项目
	project, err := s.projectRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("查找项目失败")
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}

	// 更新项目信息 - 使用充血模型方法
	project.UpdateBasicInfo(req.Name, req.Description, req.CoverImage)
	project.SetStatus(entity.ProjectStatus(req.Status))
	project.SetSchedule(req.StartDate.Time, req.Deadline.Time)
	project.SetPriorities(entity.ProjectPriority(req.Priority))

	// OwnerId check or update usually requires permission check, ignoring for now as per previous logic

	// 保存更新
	err = s.projectRepo.Update(ctx, project)
	if err != nil {
		return nil, errors.New("更新项目失败")
	}
	//先删除原来团队
	err = s.projectRepo.DeleteTeamsByProjectId(ctx, req.ID)
	if err != nil {
		return nil, errors.New("删除项目团队失败")
	}
	//再保存
	if len(req.TeamIds) > 0 {
		err = s.projectRepo.AddTeams(ctx, project.ID, req.TeamIds)
		if err != nil {
			return nil, errors.New("保存项目团队失败")
		}
	}
	return &dto.UpdateProjectResponse{
		ID:          project.ID,
		TeamIds:     req.TeamIds,
		Tags:        req.Tags,
		Priority:    int(project.Priority),
		CoverImage:  req.CoverImage,
		Deadline:    req.Deadline,
		StartDate:   req.StartDate,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		Status:      int(project.Status),
	}, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, req dto.DeleteProjectRequest) (*dto.DeleteProjectResponse, error) {
	project, err := s.projectRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("删除项目失败")
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}
	err = s.projectRepo.Delete(ctx, req.ID)
	if err != nil {
		return nil, errors.New("删除项目失败")
	}
	return &dto.DeleteProjectResponse{
		ID: project.ID,
	}, nil
}

func (s *ProjectService) GetProject(ctx context.Context, req dto.GetProjectRequest) (*dto.GetProjectResponse, error) {
	project, err := s.projectRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("获取项目失败")
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}
	teamIds, err := s.projectRepo.ListTeamIdsByProjectId(ctx, req.ID)
	if err != nil {
		return nil, errors.New("获取项目团队失败")
	}

	users, err := s.projectRepo.ListUsersByProjectId(ctx, req.ID)
	if err != nil {
		return nil, errors.New("获取项目成员失败")
	}
	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			TeamID: user.TeamID,
			Role:   user.Role,
			Status: user.Status,
		})
	}

	return &dto.GetProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerId:     project.OwnerID,
		Status:      int(project.Status),
		Deadline:    utils.NewTime(project.Deadline),
		StartDate:   utils.NewTime(project.StartDate),
		Progress:    project.Progress,
		Priority:    int(project.Priority),
		CoverImage:  project.CoverImage,
		TeamIds:     teamIds,
		Users:       userResponses,
		CreatedAt:   utils.NewTime(project.CreatedAt),
	}, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, req dto.PageRequest) (*dto.ProjectListResponse, error) {
	projects, total, err := s.projectRepo.List(ctx, req.Page, req.GetPageSize())
	if err != nil {
		return nil, errors.New("获取项目列表失败")
	}
	projectResponses := make([]*dto.ProjectResponse, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses, &dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			OwnerId:     project.OwnerID,
			Status:      int(project.Status),
			Deadline:    utils.NewTime(project.Deadline),
			StartDate:   utils.NewTime(project.StartDate),
			Progress:    project.Progress,
			Priority:    int(project.Priority),
		})
	}
	return &dto.ProjectListResponse{
		List: projectResponses,
		Page: dto.PageResponse{
			Page:     req.Page,
			PageSize: req.GetPageSize(),
			Total:    total,
		},
	}, nil
}

func (s *ProjectService) ProjectTeams(ctx context.Context) (*dto.ProjectTeamsResponse, error) {
	teams, err := s.projectRepo.ListAvailableTeams(ctx)
	if err != nil {
		return nil, errors.New("获取项目团队失败")
	}
	teamResponses := make([]*dto.TeamResponse, 0, len(teams))
	for _, team := range teams {
		teamResponses = append(teamResponses, &dto.TeamResponse{
			ID:   team.ID,
			Name: team.Name,
		})
	}
	return &dto.ProjectTeamsResponse{
		Teams: teamResponses,
	}, nil
}

func (s *ProjectService) ProjectUsers(ctx context.Context, projectID uint64) (*dto.ProjectUsersResponse, error) {
	users, err := s.projectRepo.ListUsersByProjectId(ctx, projectID)
	if err != nil {
		return nil, errors.New("获取项目成员失败")
	}
	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			TeamID: user.TeamID,
			Role:   user.Role,
			Status: user.Status,
		})
	}
	return &dto.ProjectUsersResponse{
		Users: userResponses,
	}, nil
}

func (s *ProjectService) GetProjectAvailableUsers(ctx context.Context, projectID uint64) (*dto.ProjectUsersResponse, error) {
	// 查询项目关联团队下的所有用户
	users, err := s.projectRepo.ListUsersInProjectTeams(ctx, projectID)
	if err != nil {
		return nil, errors.New("获取候选用户失败")
	}
	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			TeamID: user.TeamID,
			Role:   user.Role,
			Status: user.Status,
		})
	}
	return &dto.ProjectUsersResponse{
		Users: userResponses,
	}, nil
}

func (s *ProjectService) AddProjectUsers(ctx context.Context, req dto.AddProjectUsersRequest, projectID uint64) (*dto.ProjectUsersResponse, error) {
	// 简单校验项目是否存在
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, errors.New("查找项目失败")
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}

	// 添加用户
	err = s.projectRepo.AddUsers(ctx, projectID, req.Users)
	if err != nil {
		return nil, errors.New("添加项目成员失败")
	}

	// 返回最新的成员列表
	return s.ProjectUsers(ctx, projectID)
}

func (s *ProjectService) RemoveProjectUser(ctx context.Context, projectID uint64, userID uint64) error {
	// 简单校验项目是否存在
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return errors.New("查找项目失败")
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	// 移除用户
	err = s.projectRepo.RemoveUsers(ctx, projectID, userID)
	if err != nil {
		return errors.New("移除项目成员失败")
	}

	return nil
}
