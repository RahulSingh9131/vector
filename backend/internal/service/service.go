package service

import (
	"github.com/RahulSingh9131/vector/internal/lib/job"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
)

type Services struct {
	Auth         *AuthService
	Job          *job.JobService
	User         *UserService
	Organization *OrganizationService
	Project      *ProjectService
	Issue        *IssueService
	Label        *LabelService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	userService := NewUserService(s, repos)
	organizationService := NewOrganizationService(s, repos)
	projectService := NewProjectService(s, repos)
	issueService := NewIssueService(s, repos)
	labelService := NewLabelService(s, repos)

	return &Services{
		Job:          s.Job,
		Auth:         authService,
		User:         userService,
		Organization: organizationService,
		Project:      projectService,
		Issue:        issueService,
		Label:        labelService,
	}, nil
}

