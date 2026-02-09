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
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	userService := NewUserService(s, repos)
	organizationService := NewOrganizationService(s, repos)

	return &Services{
		Job:          s.Job,
		Auth:         authService,
		User:         userService,
		Organization: organizationService,
	}, nil
}
