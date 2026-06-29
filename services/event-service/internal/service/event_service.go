package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/event-service/internal/domain"
	"github.com/elevatecompact/spark/services/event-service/internal/events"
	"github.com/elevatecompact/spark/services/event-service/internal/repository"
)

type EventService interface {
	Create(ctx context.Context, e *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, e *domain.Event) error
	List(ctx context.Context, category string, status domain.EventStatus, page, size int) ([]domain.Event, error)

	CreateTicketTier(ctx context.Context, t *domain.EventTicketTier) error
	ListTicketTiers(ctx context.Context, eventID uuid.UUID) ([]domain.EventTicketTier, error)
	RSVP(ctx context.Context, eventID, userID uuid.UUID) error
	PurchaseTicket(ctx context.Context, tierID, userID uuid.UUID) error

	CreateSession(ctx context.Context, s *domain.EventSession) (*domain.EventSession, error)
	UpdateSession(ctx context.Context, s *domain.EventSession) error
	ListSessions(ctx context.Context, eventID uuid.UUID) ([]domain.EventSession, error)

	CreateSeries(ctx context.Context, s *domain.EventSeries) (*domain.EventSeries, error)
	GetSeries(ctx context.Context, id uuid.UUID) (*domain.EventSeries, error)
	UpdateSeries(ctx context.Context, s *domain.EventSeries) error
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	CancelEvent(ctx context.Context, id uuid.UUID) error
	ListAttendees(ctx context.Context, eventID uuid.UUID) ([]domain.EventAttendee, error)
	GetAdminStats(ctx context.Context) (*domain.EventAdminStats, error)
}

type eventService struct {
	repo     repository.EventRepository
	eventPub events.EventProducer
}

func NewEventService(repo repository.EventRepository, eventPub events.EventProducer) EventService {
	return &eventService{repo: repo, eventPub: eventPub}
}

func (s *eventService) Create(ctx context.Context, e *domain.Event) (*domain.Event, error) {
	e.ID = uuid.New()
	e.Status = domain.StatusDraft
	e.CreatedAt = time.Now().UTC()
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	s.eventPub.PublishEventCreated(ctx, e.ID)
	return e, nil
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *eventService) Update(ctx context.Context, e *domain.Event) error {
	return s.repo.Update(ctx, e)
}

func (s *eventService) List(ctx context.Context, category string, status domain.EventStatus, page, size int) ([]domain.Event, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 25
	}
	offset := (page - 1) * size
	return s.repo.List(ctx, category, status, offset, size)
}

func (s *eventService) CreateTicketTier(ctx context.Context, t *domain.EventTicketTier) error {
	t.ID = uuid.New()
	return s.repo.CreateTicketTier(ctx, t)
}

func (s *eventService) ListTicketTiers(ctx context.Context, eventID uuid.UUID) ([]domain.EventTicketTier, error) {
	return s.repo.ListTicketTiers(ctx, eventID)
}

func (s *eventService) RSVP(ctx context.Context, eventID, userID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	return s.repo.RegisterAttendee(ctx, &domain.EventAttendee{
		EventID: eventID,
		UserID:  userID,
		Status:  domain.AttendeeRegistered,
	})
}

func (s *eventService) PurchaseTicket(ctx context.Context, tierID, userID uuid.UUID) error {
	return s.repo.RegisterAttendee(ctx, &domain.EventAttendee{
		EventID:      uuid.Nil,
		TicketTierID: &tierID,
		UserID:       userID,
		Status:       domain.AttendeeRegistered,
	})
}

func (s *eventService) CreateSession(ctx context.Context, sess *domain.EventSession) (*domain.EventSession, error) {
	sess.ID = uuid.New()
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *eventService) UpdateSession(ctx context.Context, sess *domain.EventSession) error {
	return s.repo.UpdateSession(ctx, sess)
}

func (s *eventService) ListSessions(ctx context.Context, eventID uuid.UUID) ([]domain.EventSession, error) {
	return s.repo.ListSessions(ctx, eventID)
}

func (s *eventService) CreateSeries(ctx context.Context, ser *domain.EventSeries) (*domain.EventSeries, error) {
	ser.ID = uuid.New()
	ser.IsActive = true
	if err := s.repo.CreateSeries(ctx, ser); err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *eventService) GetSeries(ctx context.Context, id uuid.UUID) (*domain.EventSeries, error) {
	return s.repo.GetSeries(ctx, id)
}

func (s *eventService) UpdateSeries(ctx context.Context, ser *domain.EventSeries) error {
	return s.repo.UpdateSeries(ctx, ser)
}

func (s *eventService) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSeries(ctx, id)
}

func (s *eventService) ListAttendees(ctx context.Context, eventID uuid.UUID) ([]domain.EventAttendee, error) {
	return s.repo.ListAttendees(ctx, eventID)
}

func (s *eventService) CancelEvent(ctx context.Context, id uuid.UUID) error {
	ev, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	ev.Status = domain.StatusCancelled
	if err := s.repo.Update(ctx, ev); err != nil {
		return err
	}
	return s.eventPub.PublishEventCancelled(ctx, id)
}

func (s *eventService) GetAdminStats(ctx context.Context) (*domain.EventAdminStats, error) {
	return s.repo.GetAdminStats(ctx)
}
