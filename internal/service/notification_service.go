package service

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/2Cheetah/MedGuardianBot/internal/domain"
)

type NotificationFSMService struct {
	mu                  sync.Mutex
	sessions            map[int64]*domain.NotificationFSM // userID → FSM
	scheduleProcessor   ScheduleProcessor
	untilParser         UntilParser
	notificationService *NotificationService
}

type ScheduleProcessor interface {
	ParseSchedule(schedule string) (string, error)
}

type UntilParser interface {
	ParseText(text string) (time.Time, error)
}

func NewNotificationFSMService(sp ScheduleProcessor, up UntilParser, ns *NotificationService) *NotificationFSMService {
	return &NotificationFSMService{
		sessions:            make(map[int64]*domain.NotificationFSM),
		scheduleProcessor:   sp,
		untilParser:         up,
		notificationService: ns,
	}
}

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (nfsms *NotificationFSMService) StartSession(userID int64, chatID int64) {
	nfsms.mu.Lock()
	defer nfsms.mu.Unlock()
	nfsms.sessions[userID] = &domain.NotificationFSM{
		UserID: userID,
		ChatID: chatID,
		State:  domain.StateWaitingSchedule,
		PartialNotification: domain.Notification{
			UserID:    userID,
			ChatID:    chatID,
			CreatedAt: time.Now(),
		},
	}
}

func (nfsms *NotificationFSMService) HandleInput(userID int64, input string) (string, error) {
	nfsms.mu.Lock()
	defer nfsms.mu.Unlock()

	session, ok := nfsms.sessions[userID]
	if !ok {
		return "Please, start with /create_notification first", nil
	}

	switch session.State {
	case domain.StateWaitingSchedule:
		slog.Info("handling StateWaitingSchedule")
		schedule, err := nfsms.scheduleProcessor.ParseSchedule(input)
		if err != nil {
			return "Couldn't understand your schedule, try again?", fmt.Errorf("couldn't parse schedule string %s to a crontab. Error: %w", input, err)
		}
		session.PartialNotification.Schedule = schedule
		session.State = domain.StateWaitingUntil
		return "Until when do you want me to send notifications to you?", nil
	case domain.StateWaitingUntil:
		slog.Info("handling StateWaitingUntil")
		until, err := nfsms.untilParser.ParseText(input)
		slog.Info("parsed until", "input", input, "until", until)
		if err != nil {
			return "Couldn't understand date", fmt.Errorf("couldn't parse message %s to date. Error: %w", input, err)
		}
		// TODO: validate input
		session.PartialNotification.Until = until
		session.State = domain.StateWaitingText
		return "Notification message?", nil
	case domain.StateWaitingText:
		slog.Info("handling StateWaitingText")
		// TODO: validate input
		session.PartialNotification.Text = input
		session.PartialNotification.Status = domain.NotificationStatusActive
		if err := nfsms.notificationService.CreateNotification(&session.PartialNotification); err != nil {
			return "Something went wrong, please, try again", fmt.Errorf("couldn't create notification %v. Error: %w", session.PartialNotification, err)
		}
		session.State = domain.StateCreated
		return "Notification created!", nil
	}
	return "Something went wrong, please, try again", nil
}

func (ns *NotificationService) CreateNotification(n *domain.Notification) error {
	if n == nil {
		return fmt.Errorf("notification can't be nil")
	}
	return ns.repo.CreateNotification(n)
}
