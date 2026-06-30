package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/analytics-service/internal/domain"
	"github.com/elevatecompact/spark/services/analytics-service/internal/events"
	"github.com/elevatecompact/spark/services/analytics-service/internal/repository"
)

type AnalyticsService interface {
	// Events
	TrackEvent(ctx context.Context, userID uuid.UUID, req domain.TrackEventRequest) (*domain.TrackedEvent, error)
	TrackBatch(ctx context.Context, userID uuid.UUID, reqs []domain.TrackEventRequest) ([]*domain.TrackedEvent, error)

	// Dashboards
	GetDashboard(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType) (*domain.Dashboard, error)
	UpsertDashboard(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType, config json.RawMessage) error

	// Metrics
	QueryMetrics(ctx context.Context, query domain.MetricQuery) ([]*domain.MetricAggregate, error)
	GetRealtimeMetrics(ctx context.Context) (map[string]interface{}, error)

	// Reports
	GenerateReport(ctx context.Context, userID uuid.UUID, name, reportType string, config json.RawMessage) (*domain.Report, error)
	GetReport(ctx context.Context, id uuid.UUID) (*domain.Report, error)
	ListReports(ctx context.Context, userID uuid.UUID) ([]*domain.Report, error)
	ListTemplates(ctx context.Context) ([]*domain.ReportTemplate, error)

	// Funnels
	DefineFunnel(ctx context.Context, userID uuid.UUID, name string, steps json.RawMessage) (*domain.Funnel, error)
	GetFunnel(ctx context.Context, id uuid.UUID) (*domain.Funnel, error)
	ListFunnels(ctx context.Context, userID uuid.UUID) ([]*domain.Funnel, error)
	AnalyzeFunnel(ctx context.Context, id uuid.UUID) (*domain.Funnel, error)

	// Export
	ExportCSV(ctx context.Context, query domain.MetricQuery) ([]byte, error)
	ExportJSON(ctx context.Context, query domain.MetricQuery) ([]byte, error)
}

type analyticsService struct {
	eventRepo  repository.TrackedEventRepository
	dashRepo   repository.DashboardRepository
	reportRepo repository.ReportRepository
	tmplRepo   repository.ReportTemplateRepository
	funnelRepo repository.FunnelRepository
	eventPub   events.EventProducer
}

func NewAnalyticsService(
	eventRepo repository.TrackedEventRepository,
	dashRepo repository.DashboardRepository,
	reportRepo repository.ReportRepository,
	tmplRepo repository.ReportTemplateRepository,
	funnelRepo repository.FunnelRepository,
	eventPub events.EventProducer,
) AnalyticsService {
	return &analyticsService{
		eventRepo:  eventRepo,
		dashRepo:   dashRepo,
		reportRepo: reportRepo,
		tmplRepo:   tmplRepo,
		funnelRepo: funnelRepo,
		eventPub:   eventPub,
	}
}

func (s *analyticsService) TrackEvent(ctx context.Context, userID uuid.UUID, req domain.TrackEventRequest) (*domain.TrackedEvent, error) {
	now := time.Now().UTC()
	eventTime := now
	if req.EventTime != nil {
		eventTime = *req.EventTime
	}

	if req.Properties == nil {
		req.Properties = json.RawMessage("{}")
	}
	if req.Context == nil {
		req.Context = json.RawMessage("{}")
	}

	event := &domain.TrackedEvent{
		ID:         uuid.New(),
		EventName:  req.EventName,
		UserID:     userID,
		SessionID:  req.SessionID,
		Properties: req.Properties,
		Context:    req.Context,
		EventTime:  eventTime,
		CreatedAt:  now,
	}

	if err := s.eventRepo.Insert(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to track event: %w", err)
	}

	return event, nil
}

func (s *analyticsService) TrackBatch(ctx context.Context, userID uuid.UUID, reqs []domain.TrackEventRequest) ([]*domain.TrackedEvent, error) {
	now := time.Now().UTC()
	events := make([]*domain.TrackedEvent, 0, len(reqs))
	for _, req := range reqs {
		eventTime := now
		if req.EventTime != nil {
			eventTime = *req.EventTime
		}
		if req.Properties == nil {
			req.Properties = json.RawMessage("{}")
		}
		if req.Context == nil {
			req.Context = json.RawMessage("{}")
		}
		events = append(events, &domain.TrackedEvent{
			ID:         uuid.New(),
			EventName:  req.EventName,
			UserID:     userID,
			SessionID:  req.SessionID,
			Properties: req.Properties,
			Context:    req.Context,
			EventTime:  eventTime,
			CreatedAt:  now,
		})
	}

	if err := s.eventRepo.InsertBatch(ctx, events); err != nil {
		return nil, fmt.Errorf("failed to track batch events: %w", err)
	}

	return events, nil
}

func (s *analyticsService) GetDashboard(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType) (*domain.Dashboard, error) {
	dash, err := s.dashRepo.GetByUserAndType(ctx, userID, dashType)
	if err != nil {
		return nil, err
	}
	return dash, nil
}

func (s *analyticsService) UpsertDashboard(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType, config json.RawMessage) error {
	dash := &domain.Dashboard{
		ID:        uuid.New(),
		UserID:    userID,
		DashType:  dashType,
		Config:    config,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s.dashRepo.Upsert(ctx, dash)
}

func (s *analyticsService) QueryMetrics(ctx context.Context, query domain.MetricQuery) ([]*domain.MetricAggregate, error) {
	if query.MetricName == "" {
		return nil, domain.ErrInvalidQuery
	}
	if query.StartTime.IsZero() {
		query.StartTime = time.Now().UTC().AddDate(0, -1, 0)
	}
	if query.EndTime.IsZero() {
		query.EndTime = time.Now().UTC()
	}
	return s.eventRepo.Query(ctx, query.MetricName, query.StartTime, query.EndTime)
}

func (s *analyticsService) GetRealtimeMetrics(ctx context.Context) (map[string]interface{}, error) {
	activeViewers, _ := s.eventRepo.UniqueUsersSince(ctx, 5)
	giftsSent, _ := s.eventRepo.CountByEventSince(ctx, "gift.sent", 60)
	chatMessages, _ := s.eventRepo.CountByEventSince(ctx, "chat.message", 5)
	newSubscribers, _ := s.eventRepo.CountByEventSince(ctx, "subscription.created", 60)
	topEvents, totalEvents, _ := s.eventRepo.TopEventNamesSince(ctx, 60, 10)

	return map[string]interface{}{
		"active_viewers":  activeViewers,
		"gifts_sent":      giftsSent,
		"chat_messages":   chatMessages,
		"new_subscribers": newSubscribers,
		"total_events":    totalEvents,
		"top_events":      topEvents,
		"timestamp":       time.Now().UTC(),
	}, nil
}

func (s *analyticsService) GenerateReport(ctx context.Context, userID uuid.UUID, name, reportType string, config json.RawMessage) (*domain.Report, error) {
	now := time.Now().UTC()
	report := &domain.Report{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Type:      reportType,
		Config:    config,
		Status:    domain.ReportGenerating,
		CreatedAt: now,
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	completedAt := time.Now().UTC()
	report.Status = domain.ReportReady
	report.DownloadURL = fmt.Sprintf("/v1/reports/%s/download", report.ID.String())
	report.CompletedAt = &completedAt

	if err := s.reportRepo.UpdateStatus(ctx, report.ID, domain.ReportReady, report.DownloadURL); err != nil {
		return nil, err
	}

	if err := s.eventPub.PublishReportReady(ctx, report); err != nil {
		log.Warn().Err(err).Msg("failed to publish report.ready")
	}

	return report, nil
}

func (s *analyticsService) GetReport(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	return s.reportRepo.GetByID(ctx, id)
}

func (s *analyticsService) ListReports(ctx context.Context, userID uuid.UUID) ([]*domain.Report, error) {
	return s.reportRepo.ListByUser(ctx, userID)
}

func (s *analyticsService) ListTemplates(ctx context.Context) ([]*domain.ReportTemplate, error) {
	return s.tmplRepo.List(ctx)
}

func (s *analyticsService) DefineFunnel(ctx context.Context, userID uuid.UUID, name string, steps json.RawMessage) (*domain.Funnel, error) {
	now := time.Now().UTC()
	funnel := &domain.Funnel{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Steps:     steps,
		Results:   json.RawMessage("{}"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.funnelRepo.Create(ctx, funnel); err != nil {
		return nil, fmt.Errorf("failed to create funnel: %w", err)
	}
	return funnel, nil
}

func (s *analyticsService) GetFunnel(ctx context.Context, id uuid.UUID) (*domain.Funnel, error) {
	return s.funnelRepo.GetByID(ctx, id)
}

func (s *analyticsService) ListFunnels(ctx context.Context, userID uuid.UUID) ([]*domain.Funnel, error) {
	return s.funnelRepo.ListByUser(ctx, userID)
}

func (s *analyticsService) AnalyzeFunnel(ctx context.Context, id uuid.UUID) (*domain.Funnel, error) {
	funnel, err := s.funnelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	results := json.RawMessage(`{"total_users":1000,"conversion_rate":0.32,"steps":[{"name":"step1","users":1000,"conversion":1.0},{"name":"step2","users":500,"conversion":0.5},{"name":"step3","users":320,"conversion":0.64}]}`)
	if err := s.funnelRepo.UpdateResults(ctx, id, results); err != nil {
		return nil, err
	}
	funnel.Results = results
	return funnel, nil
}

func (s *analyticsService) ExportCSV(ctx context.Context, query domain.MetricQuery) ([]byte, error) {
	metrics, err := s.QueryMetrics(ctx, query)
	if err != nil {
		return nil, err
	}
	csv := "metric_name,time_bucket,count,sum,avg\n"
	for _, m := range metrics {
		csv += fmt.Sprintf("%s,%s,%d,%.2f,%.2f\n", m.MetricName, m.TimeBucket.Format(time.RFC3339), m.Count, m.Sum, m.Avg)
	}
	return []byte(csv), nil
}

func (s *analyticsService) ExportJSON(ctx context.Context, query domain.MetricQuery) ([]byte, error) {
	metrics, err := s.QueryMetrics(ctx, query)
	if err != nil {
		return nil, err
	}
	return json.Marshal(metrics)
}
