package entity

import "time"

type ProjectStatus int

const (
	ProjectStatusActive    ProjectStatus = 1
	ProjectStatusCompleted ProjectStatus = 2
	ProjectStatusArchived  ProjectStatus = 3
)

type ProjectPriority int

const (
	ProjectPriorityP0 ProjectPriority = 0
	ProjectPriorityP1 ProjectPriority = 1
	ProjectPriorityP2 ProjectPriority = 2
	ProjectPriorityP3 ProjectPriority = 3
)

type Project struct {
	BaseEntity
	Name        string
	Description string
	OwnerID     uint64
	Status      ProjectStatus
	Deadline    time.Time
	StartDate   time.Time
	Progress    int
	Priority    ProjectPriority
	CoverImage  string
}

// NewProject 创建新项目
func NewProject(name, description string, ownerID uint64) *Project {
	return &Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Status:      ProjectStatusActive,
		Progress:    0,
		Priority:    ProjectPriorityP2,
	}
}

// UpdateBasicInfo 更新基本信息
func (p *Project) UpdateBasicInfo(name, description string, coverImage string) {
	p.Name = name
	p.Description = description
	p.CoverImage = coverImage
}

// SetStatus 设置状态
func (p *Project) SetStatus(status ProjectStatus) {
	p.Status = status
}

// SetSchedule 设置进度安排
func (p *Project) SetSchedule(startDate, deadline time.Time) {
	p.StartDate = startDate
	p.Deadline = deadline
}

// SetPriorities 设置优先级和标签
func (p *Project) SetPriorities(priority ProjectPriority) {
	p.Priority = priority
}

func (p *Project) IsActive() bool {
	return p.Status == ProjectStatusActive && !p.IsDeleted()
}
