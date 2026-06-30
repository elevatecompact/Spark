package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/notification-service/internal/domain"
	"github.com/elevatecompact/spark/services/notification-service/internal/events"
	"github.com/elevatecompact/spark/services/notification-service/internal/processor"
	"github.com/elevatecompact/spark/services/notification-service/internal/repository"
)

type NotificationService interface {
	// Inbox
	ListNotifications(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Notification, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id, userID uuid.UUID) error

	// Preferences
	GetPreferences(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreference, error)
	UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs string) error

	// Send
	SendNotification(ctx context.Context, req domain.SendNotificationRequest) (*domain.Notification, error)
	SendBatch(ctx context.Context, reqs []domain.SendNotificationRequest) ([]*domain.Notification, error)

	// Devices
	RegisterDevice(ctx context.Context, userID uuid.UUID, req domain.RegisterDeviceRequest) (*domain.PushDevice, error)
	UnregisterDevice(ctx context.Context, id, userID uuid.UUID) error
	ListDevices(ctx context.Context, userID uuid.UUID) ([]*domain.PushDevice, error)

	// Templates
	ListTemplates(ctx context.Context) ([]*domain.Template, error)
	CreateTemplate(ctx context.Context, t *domain.Template) error
	UpdateTemplate(ctx context.Context, t *domain.Template) error

	// Admin
	TestPush(ctx context.Context, userID uuid.UUID) error
	TestEmail(ctx context.Context, email string) error
	DeliveryStats(ctx context.Context) (map[string]interface{}, error)
}

type notifService struct {
	notifRepo repository.NotificationRepository
	prefRepo  repository.PreferenceRepository
	devRepo   repository.DeviceRepository
	tmplRepo  repository.TemplateRepository
	push      processor.PushProcessor
	email     processor.EmailProcessor
	sms       processor.SMSProcessor
	eventPub  events.EventProducer
	pushOn    bool
	emailOn   bool
	smsOn     bool
	inappOn   bool
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	prefRepo repository.PreferenceRepository,
	devRepo repository.DeviceRepository,
	tmplRepo repository.TemplateRepository,
	push processor.PushProcessor,
	email processor.EmailProcessor,
	sms processor.SMSProcessor,
	eventPub events.EventProducer,
	pushOn, emailOn, smsOn, inappOn bool,
) NotificationService {
	return &notifService{
		notifRepo: notifRepo,
		prefRepo:  prefRepo,
		devRepo:   devRepo,
		tmplRepo:  tmplRepo,
		push:      push,
		email:     email,
		sms:       sms,
		eventPub:  eventPub,
		pushOn:    pushOn,
		emailOn:   emailOn,
		smsOn:     smsOn,
		inappOn:   inappOn,
	}
}

func (s *notifService) ListNotifications(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Notification, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.notifRepo.ListByUser(ctx, userID, cursor, limit)
}

func (s *notifService) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	return s.notifRepo.MarkRead(ctx, id)
}

func (s *notifService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.notifRepo.MarkAllRead(ctx, userID)
}

func (s *notifService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.notifRepo.Delete(ctx, id)
}

func (s *notifService) GetPreferences(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreference, error) {
	pref, err := s.prefRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		pref = &domain.NotificationPreference{
			UserID:      userID,
			Preferences: `{}`,
		}
	}
	return pref, nil
}

func (s *notifService) UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs string) error {
	return s.prefRepo.Upsert(ctx, &domain.NotificationPreference{
		UserID:      userID,
		Preferences: prefs,
	})
}

func (s *notifService) SendNotification(ctx context.Context, req domain.SendNotificationRequest) (*domain.Notification, error) {
	now := time.Now().UTC()
	notif := &domain.Notification{
		ID:        uuid.New(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Data:      req.Data,
		Channel:   domain.ChannelInApp,
		CreatedAt: now,
	}

	if s.inappOn {
		if err := s.notifRepo.Insert(ctx, notif); err != nil {
			return nil, fmt.Errorf("failed to insert notification: %w", err)
		}
	}

	if s.pushOn {
		devices, err := s.devRepo.ListByUser(ctx, req.UserID)
		if err == nil {
			for _, d := range devices {
				if err := s.push.Send(ctx, d.Token, req.Title, req.Body, nil); err != nil {
					log.Warn().Err(err).Str("device", d.ID.String()).Msg("push send failed")
				}
			}
		}
	}

	if s.emailOn {
		if err := s.email.Send(ctx, req.UserID.String()+"@test.com", req.Title, req.Body); err != nil {
			log.Warn().Err(err).Msg("email send failed")
		}
	}

	if err := s.eventPub.PublishDelivered(ctx, notif.ID, req.UserID, string(domain.ChannelInApp), "delivered"); err != nil {
		log.Warn().Err(err).Msg("failed to publish delivery event")
	}

	return notif, nil
}

func (s *notifService) SendBatch(ctx context.Context, reqs []domain.SendNotificationRequest) ([]*domain.Notification, error) {
	notifs := make([]*domain.Notification, 0, len(reqs))
	for _, req := range reqs {
		n, err := s.SendNotification(ctx, req)
		if err != nil {
			return notifs, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (s *notifService) RegisterDevice(ctx context.Context, userID uuid.UUID, req domain.RegisterDeviceRequest) (*domain.PushDevice, error) {
	dev := &domain.PushDevice{
		ID:        uuid.New(),
		UserID:    userID,
		Platform:  req.Platform,
		Token:     req.Token,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.devRepo.Register(ctx, dev); err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}
	return dev, nil
}

func (s *notifService) UnregisterDevice(ctx context.Context, id, userID uuid.UUID) error {
	return s.devRepo.Deactivate(ctx, id)
}

func (s *notifService) ListDevices(ctx context.Context, userID uuid.UUID) ([]*domain.PushDevice, error) {
	return s.devRepo.ListByUser(ctx, userID)
}

func (s *notifService) ListTemplates(ctx context.Context) ([]*domain.Template, error) {
	return s.tmplRepo.List(ctx)
}

func (s *notifService) CreateTemplate(ctx context.Context, t *domain.Template) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().UTC()
	return s.tmplRepo.Create(ctx, t)
}

func (s *notifService) UpdateTemplate(ctx context.Context, t *domain.Template) error {
	return s.tmplRepo.Update(ctx, t)
}

func (s *notifService) TestPush(ctx context.Context, userID uuid.UUID) error {
	devices, err := s.devRepo.ListByUser(ctx, userID)
	if err != nil {
		return err
	}
	for _, d := range devices {
		if err := s.push.Send(ctx, d.Token, "Test Notification", "This is a test push notification", nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *notifService) TestEmail(ctx context.Context, email string) error {
	return s.email.Send(ctx, email, "Test Email", "This is a test email from the notification service.")
}

func (s *notifService) DeliveryStats(ctx context.Context) (map[string]interface{}, error) {
	stats, err := s.notifRepo.DeliveryStats(ctx)
	if err != nil {
		return nil, err
	}
	total := stats.PushDelivered + stats.EmailDelivered + stats.SMSDelivered + stats.InAppDelivered
	bounceRate := 0.0
	if total > 0 {
		bounceRate = float64(stats.EmailDelivered) / float64(total)
	}
	return map[string]interface{}{
		"push_delivered":  stats.PushDelivered,
		"email_delivered": stats.EmailDelivered,
		"sms_delivered":   stats.SMSDelivered,
		"inapp_delivered": stats.InAppDelivered,
		"bounce_rate":     bounceRate,
	}, nil
}
