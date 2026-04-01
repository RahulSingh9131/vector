package service

import (
	"context"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

type NotificationService struct {
	server           *server.Server
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(s *server.Server, repos *repository.Repositories) *NotificationService {
	return &NotificationService{
		server:           s,
		notificationRepo: repos.Notification,
	}
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID uuid.UUID, filters models.NotificationFilters) ([]models.Notification, error) {
	return s.notificationRepo.ListByUser(ctx, userID, filters)
}

func (s *NotificationService) GetNotification(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}
