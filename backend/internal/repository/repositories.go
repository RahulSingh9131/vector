package repository

import "github.com/RahulSingh9131/vector/internal/server"

type Repositories struct {
	User               *UserRepository
	Organization       *OrganizationRepository
	OrganizationMember *OrganizationMemberRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User:               NewUserRepository(s),
		Organization:       NewOrganizationRepository(s),
		OrganizationMember: NewOrganizationMemberRepository(s),
	}
}
