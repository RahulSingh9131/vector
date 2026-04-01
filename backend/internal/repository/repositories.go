package repository

import "github.com/RahulSingh9131/vector/internal/server"

type Repositories struct {
	User               *UserRepository
	Organization       *OrganizationRepository
	OrganizationMember *OrganizationMemberRepository
	Project            *ProjectRepository
	ProjectMember      *ProjectMemberRepository
	Issue              *IssueRepository
	Label              *LabelRepository
	Comment            *CommentRepository
	Activity           *ActivityRepository
	Notification       *NotificationRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User:               NewUserRepository(s),
		Organization:       NewOrganizationRepository(s),
		OrganizationMember: NewOrganizationMemberRepository(s),
		Project:            NewProjectRepository(s),
		ProjectMember:      NewProjectMemberRepository(s),
		Issue:              NewIssueRepository(s),
		Label:              NewLabelRepository(s),
		Comment:            NewCommentRepository(s),
		Activity:           NewActivityRepository(s),
		Notification:       NewNotificationRepository(s, s.DB.Pool),
	}
}



