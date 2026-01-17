package repository

import (
	"FLOWGO/internal/domain/entity"
	"context"
)

type ProjectsRepository interface {
	BaseRepository[entity.Project]
	ListAvailableTeams(ctx context.Context) ([]*entity.Team, error)
	DeleteTeamsByProjectId(ctx context.Context, projectId uint64) error
	AddTeams(ctx context.Context, projectId uint64, teamIds []uint64) error
	ListTeamIdsByProjectId(ctx context.Context, projectId uint64) ([]uint64, error)
	ListUsersByProjectId(ctx context.Context, projectId uint64) ([]*entity.User, error)
	ListAvailableUsers(ctx context.Context) ([]*entity.User, error)
	AddUsers(ctx context.Context, projectId uint64, userIds []uint64) error
	RemoveUsers(ctx context.Context, projectId uint64, userId uint64) error
	ListUsersInProjectTeams(ctx context.Context, projectId uint64) ([]*entity.User, error)
}
