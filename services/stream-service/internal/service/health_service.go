package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type HealthService struct {
	mu      sync.RWMutex
	reports map[uuid.UUID]*StreamHealthTracker
}

type StreamHealthTracker struct {
	StreamID    uuid.UUID
	Current     domain.StreamHealth
	History     []domain.HealthReport
	LastUpdated time.Time
	maxHistory  int
}

func NewHealthService() *HealthService {
	return &HealthService{
		reports: make(map[uuid.UUID]*StreamHealthTracker),
	}
}

func (s *HealthService) ReportHealth(ctx context.Context, streamID uuid.UUID, report domain.HealthReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tracker, ok := s.reports[streamID]
	if !ok {
		tracker = &StreamHealthTracker{
			StreamID:   streamID,
			History:    make([]domain.HealthReport, 0, 100),
			maxHistory: 1000,
		}
		s.reports[streamID] = tracker
	}

	tracker.History = append(tracker.History, report)
	if len(tracker.History) > tracker.maxHistory {
		tracker.History = tracker.History[len(tracker.History)-tracker.maxHistory:]
	}

	tracker.Current = calculateStreamHealth(report)
	tracker.LastUpdated = time.Now()

	return nil
}

func (s *HealthService) GetHealth(ctx context.Context, streamID uuid.UUID) (*domain.StreamHealth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tracker, ok := s.reports[streamID]
	if !ok {
		return &domain.StreamHealth{
			StreamID:        streamID,
			ConnectionScore: 100,
			Status:          "unknown",
		}, nil
	}

	return &tracker.Current, nil
}

func (s *HealthService) GetHealthHistory(ctx context.Context, streamID uuid.UUID, duration time.Duration) ([]domain.HealthReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tracker, ok := s.reports[streamID]
	if !ok {
		return nil, nil
	}

	cutoff := time.Now().Add(-duration)
	var filtered []domain.HealthReport
	for _, r := range tracker.History {
		if r.Timestamp.After(cutoff) {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

func (s *HealthService) DetectAnomalies(ctx context.Context, streamID uuid.UUID) (*domain.AnomalyReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tracker, ok := s.reports[streamID]
	if !ok {
		return &domain.AnomalyReport{
			StreamID:  streamID,
			Anomalies: []domain.Anomaly{},
			Healthy:   true,
		}, nil
	}

	report := &domain.AnomalyReport{
		StreamID:  streamID,
		Anomalies: []domain.Anomaly{},
		Healthy:   true,
	}

	if len(tracker.History) < 5 {
		return report, nil
	}

	recent := tracker.History
	if len(recent) > 30 {
		recent = recent[len(recent)-30:]
	}

	var totalBitrate int
	var totalFPS float64
	var totalPacketLoss float64
	for _, r := range recent {
		totalBitrate += r.Bitrate
		totalFPS += r.FPS
		totalPacketLoss += r.PacketLoss
	}
	avgBitrate := totalBitrate / len(recent)
	avgFPS := totalFPS / float64(len(recent))
	avgPacketLoss := totalPacketLoss / float64(len(recent))

	if avgPacketLoss > 5.0 {
		report.Anomalies = append(report.Anomalies, domain.Anomaly{
			Type:        "high_packet_loss",
			Severity:    "warning",
			Value:       avgPacketLoss,
			Threshold:   5.0,
			Description: fmt.Sprintf("Average packet loss %.1f%% exceeds threshold of 5%%", avgPacketLoss),
		})
		report.Healthy = false
	}

	last := recent[len(recent)-1]
	if avgBitrate > 0 && last.Bitrate < int(float64(avgBitrate)*0.5) {
		report.Anomalies = append(report.Anomalies, domain.Anomaly{
			Type:        "bitrate_drop",
			Severity:    "critical",
			Value:       float64(last.Bitrate),
			Threshold:   float64(avgBitrate) * 0.5,
			Description: fmt.Sprintf("Bitrate dropped from avg %d to %d", avgBitrate, last.Bitrate),
		})
		report.Healthy = false
	}

	if last.FPS > 0 && avgFPS > 0 && last.FPS < avgFPS*0.5 {
		report.Anomalies = append(report.Anomalies, domain.Anomaly{
			Type:        "fps_drop",
			Severity:    "warning",
			Value:       last.FPS,
			Threshold:   avgFPS * 0.5,
			Description: fmt.Sprintf("FPS dropped from avg %.1f to %.1f", avgFPS, last.FPS),
		})
		report.Healthy = false
	}

	if last.RoundTripTime > 500 {
		report.Anomalies = append(report.Anomalies, domain.Anomaly{
			Type:        "high_latency",
			Severity:    "warning",
			Value:       last.RoundTripTime,
			Threshold:   500,
			Description: fmt.Sprintf("Round trip time %.0fms exceeds 500ms threshold", last.RoundTripTime),
		})
	}

	if last.Jitter > 100 {
		report.Anomalies = append(report.Anomalies, domain.Anomaly{
			Type:        "high_jitter",
			Severity:    "info",
			Value:       last.Jitter,
			Threshold:   100,
			Description: fmt.Sprintf("Jitter %.0fms exceeds 100ms threshold", last.Jitter),
		})
	}

	return report, nil
}

func calculateStreamHealth(report domain.HealthReport) domain.StreamHealth {
	score := 100.0

	if report.PacketLoss > 0 {
		score -= math.Min(report.PacketLoss*2, 40)
	}

	if report.RoundTripTime > 100 {
		excess := (report.RoundTripTime - 100) / 10
		score -= math.Min(excess, 20)
	}

	if report.Jitter > 30 {
		excess := (report.Jitter - 30) / 5
		score -= math.Min(excess, 20)
	}

	if score < 0 {
		score = 0
	}

	status := "good"
	if score < 80 {
		status = "degraded"
	}
	if score < 50 {
		status = "poor"
	}
	if score < 20 {
		status = "critical"
	}

	return domain.StreamHealth{
		StreamID:        report.StreamID,
		Bitrate:         report.Bitrate,
		FPS:             report.FPS,
		PacketLoss:      report.PacketLoss,
		RoundTripTime:   report.RoundTripTime,
		Jitter:          report.Jitter,
		ConnectionScore: int(math.Round(score)),
		Status:          status,
		LastUpdated:     time.Now(),
	}
}

func (s *HealthService) RemoveStream(ctx context.Context, streamID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.reports, streamID)
}
