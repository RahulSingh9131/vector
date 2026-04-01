package service

import (
	"github.com/RahulSingh9131/vector/internal/events"
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
	Comment      *CommentService
	Activity     *ActivityService
	Notification *NotificationService
	Search       *SearchService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	// Create event bus and publisher
	eventBus := events.NewAsynqEventBus(s.Job.Client, s.Logger)
	eventPublisher := events.NewEventPublisher(eventBus, s.Logger)

	authService := NewAuthService(s)
	userService := NewUserService(s, repos)
	organizationService := NewOrganizationService(s, repos)
	projectService := NewProjectService(s, repos, eventPublisher)
	issueService := NewIssueService(s, repos, eventPublisher)
	labelService := NewLabelService(s, repos, eventPublisher)
	commentService := NewCommentService(s, repos, eventPublisher)
	activityService := NewActivityService(s, repos)
	notificationService := NewNotificationService(s, repos)
	searchService := NewSearchService(s, repos)

	return &Services{
		Job:          s.Job,
		Auth:         authService,
		User:         userService,
		Organization: organizationService,
		Project:      projectService,
		Issue:        issueService,
		Label:        labelService,
		Comment:      commentService,
		Activity:     activityService,
		Notification: notificationService,
		Search:       searchService,
	}, nil
}
