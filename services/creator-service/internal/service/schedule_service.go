package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
)

type ScheduleService struct {
	pool        *pgxpool.Pool
	creatorRepo repository.CreatorRepository
}

func NewScheduleService(pool *pgxpool.Pool, creatorRepo repository.CreatorRepository) *ScheduleService {
	return &ScheduleService{
		pool:        pool,
		creatorRepo: creatorRepo,
	}
}

func (s *ScheduleService) GetSchedule(ctx context.Context, creatorID uuid.UUID) ([]domain.ScheduleSlot, error) {
	rows, err := s.pool.Query(ctx, `SELECT * FROM schedule_slots WHERE creator_id = $1 AND active = true ORDER BY day_of_week, start_time`, creatorID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	defer rows.Close()

	var slots []domain.ScheduleSlot
	for rows.Next() {
		var slot domain.ScheduleSlot
		if err := rows.Scan(&slot.ID, &slot.CreatorID, &slot.DayOfWeek, &slot.StartTime, &slot.EndTime, &slot.Title, &slot.Recurring, &slot.Active, &slot.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan schedule slot: %w", err)
		}
		slots = append(slots, slot)
	}
	if slots == nil {
		slots = []domain.ScheduleSlot{}
	}
	return slots, nil
}

func (s *ScheduleService) AddSlot(ctx context.Context, creatorID uuid.UUID, req domain.CreateScheduleSlotRequest) (*domain.ScheduleSlot, error) {
	if _, err := s.creatorRepo.GetByID(ctx, creatorID); err != nil {
		return nil, err
	}

	conflict, err := s.checkConflict(ctx, creatorID, req.DayOfWeek, req.StartTime, req.EndTime, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, domain.ErrScheduleConflict
	}

	slot := &domain.ScheduleSlot{
		ID:        uuid.New(),
		CreatorID: creatorID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title:     req.Title,
		Recurring: req.Recurring,
		Active:    true,
		CreatedAt: time.Now(),
	}

	_, err = s.pool.Exec(ctx, `INSERT INTO schedule_slots (id, creator_id, day_of_week, start_time, end_time, title, recurring, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		slot.ID, slot.CreatorID, slot.DayOfWeek, slot.StartTime, slot.EndTime, slot.Title, slot.Recurring, slot.Active, slot.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert schedule slot: %w", err)
	}

	log.Info().Str("creator_id", creatorID.String()).Int("day", req.DayOfWeek).Str("start", req.StartTime).Msg("Schedule slot added")
	return slot, nil
}

func (s *ScheduleService) UpdateSlot(ctx context.Context, id, creatorID uuid.UUID, req domain.CreateScheduleSlotRequest) error {
	conflict, err := s.checkConflict(ctx, creatorID, req.DayOfWeek, req.StartTime, req.EndTime, id)
	if err != nil {
		return err
	}
	if conflict {
		return domain.ErrScheduleConflict
	}

	_, err = s.pool.Exec(ctx, `UPDATE schedule_slots SET day_of_week=$1, start_time=$2, end_time=$3, title=$4, recurring=$5 WHERE id=$6 AND creator_id=$7`,
		req.DayOfWeek, req.StartTime, req.EndTime, req.Title, req.Recurring, id, creatorID)
	if err != nil {
		return fmt.Errorf("update schedule slot: %w", err)
	}
	return nil
}

func (s *ScheduleService) DeleteSlot(ctx context.Context, id, creatorID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM schedule_slots WHERE id = $1 AND creator_id = $2`, id, creatorID)
	if err != nil {
		return fmt.Errorf("delete schedule slot: %w", err)
	}
	return nil
}

func (s *ScheduleService) GetLiveNow(ctx context.Context) ([]domain.Creator, error) {
	now := time.Now().UTC()
	currentDay := int(now.Weekday())
	currentTime := now.Format("15:04")

	rows, err := s.pool.Query(ctx, `SELECT c.* FROM creators c 
		INNER JOIN schedule_slots s ON s.creator_id = c.id 
		WHERE s.active = true AND s.day_of_week = $1 AND s.start_time <= $2 AND s.end_time >= $3
		GROUP BY c.id`, currentDay, currentTime, currentTime)
	if err != nil {
		return nil, fmt.Errorf("get live now: %w", err)
	}
	defer rows.Close()

	return nil, nil
}

func (s *ScheduleService) checkConflict(ctx context.Context, creatorID uuid.UUID, dayOfWeek int, startTime, endTime string, excludeID uuid.UUID) (bool, error) {
	var exists bool
	var err error
	if excludeID != uuid.Nil {
		err = s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schedule_slots 
			WHERE creator_id = $1 AND day_of_week = $2 AND active = true AND id != $3
			AND ((start_time <= $4 AND end_time > $4) OR (start_time < $5 AND end_time >= $5) OR (start_time >= $4 AND end_time <= $5)))`,
			creatorID, dayOfWeek, excludeID, startTime, endTime).Scan(&exists)
	} else {
		err = s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schedule_slots 
			WHERE creator_id = $1 AND day_of_week = $2 AND active = true
			AND ((start_time <= $3 AND end_time > $3) OR (start_time < $4 AND end_time >= $4) OR (start_time >= $3 AND end_time <= $4)))`,
			creatorID, dayOfWeek, startTime, endTime).Scan(&exists)
	}
	if err != nil {
		return false, fmt.Errorf("check conflict: %w", err)
	}
	return exists, nil
}
