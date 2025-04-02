package service

import (
	"fmt"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (ns *NotificationService) CreateNotification(n *domain.Notification) error {
	if n == nil {
		return fmt.Errorf("notification can't be nil")
	}
	return ns.repo.CreateNotification(n)
}
