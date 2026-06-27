package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/events"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
)

type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationApproved VerificationStatus = "approved"
	VerificationRejected VerificationStatus = "rejected"
	VerificationNone     VerificationStatus = "none"
)

type VerificationService struct {
	creatorRepo repository.CreatorRepository
	producer    *events.Producer
}

func NewVerificationService(creatorRepo repository.CreatorRepository, producer *events.Producer) *VerificationService {
	return &VerificationService{
		creatorRepo: creatorRepo,
		producer:    producer,
	}
}

type VerificationRequest struct {
	CreatorID uuid.UUID json:"creator_id"
	Status    VerificationStatus json:"status"
	AdminID   uuid.UUID          json:"admin_id,omitempty"
	Reason    string             json:"reason,omitempty"
	Documents []string           json:"documents,omitempty"
	CreatedAt time.Time          json:"created_at"
	UpdatedAt time.Time          json:"updated_at"
}

func (s *VerificationService) RequestVerification(ctx context.Context, creatorID uuid.UUID, documents []string) error {
	creator, err := s.creatorRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}

	if creator.Verified {
		return domain.ErrAlreadyVerified
	}

	existing := s.getPendingRequest(ctx, creatorID)
	if existing != nil {
		return domain.ErrVerificationPending
	}

	req := &VerificationRequest{
		CreatorID: creatorID,
		Status:    VerificationPending,
		Documents: documents,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.saveVerificationRequest(ctx, req); err != nil {
		return fmt.Errorf("save verification request: %w", err)
	}

	log.Info().Str("creator_id", creatorID.String()).Msg("Verification requested")
	return nil
}

func (s *VerificationService) ApproveVerification(ctx context.Context, creatorID, adminID uuid.UUID) error {
	creator, err := s.creatorRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}

	if creator.Verified {
		return domain.ErrAlreadyVerified
	}

	now := time.Now()
	creator.Verified = true
	creator.VerifiedAt = &now
	creator.UpdatedAt = now

	if err := s.creatorRepo.Update(ctx, creator); err != nil {
		return fmt.Errorf("update creator verification: %w", err)
	}

	if err := s.deletePendingRequest(ctx, creatorID); err != nil {
		log.Warn().Err(err).Msg("failed to delete pending verification request")
	}

	if err := s.producer.CreatorVerified(ctx, creatorID.String(), adminID.String()); err != nil {
		log.Warn().Err(err).Msg("failed to emit CreatorVerified event")
	}

	log.Info().Str("creator_id", creatorID.String()).Str("admin_id", adminID.String()).Msg("Creator verified")
	return nil
}

func (s *VerificationService) RejectVerification(ctx context.Context, creatorID uuid.UUID, reason string) error {
	creator, err := s.creatorRepo.GetByID(ctx, creatorID)
	if err != nil {
		return err
	}
	if creator.Verified {
		return domain.ErrAlreadyVerified
	}

	req := s.getPendingRequest(ctx, creatorID)
	if req == nil {
		return domain.ErrVerificationPending
	}

	req.Status = VerificationRejected
	req.Reason = reason
	req.UpdatedAt = time.Now()

	if err := s.saveVerificationRequest(ctx, req); err != nil {
		return fmt.Errorf("save rejected verification: %w", err)
	}

	log.Info().Str("creator_id", creatorID.String()).Str("reason", reason).Msg("Verification rejected")
	return nil
}

func (s *VerificationService) GetVerificationStatus(ctx context.Context, creatorID uuid.UUID) (*VerificationRequest, error) {
	creator, err := s.creatorRepo.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	if creator.Verified {
		return &VerificationRequest{
			CreatorID: creatorID,
			Status:    VerificationApproved,
			UpdatedAt: *creator.VerifiedAt,
		}, nil
	}

	req := s.getPendingRequest(ctx, creatorID)
	if req != nil {
		return req, nil
	}

	return &VerificationRequest{
		CreatorID: creatorID,
		Status:    VerificationNone,
	}, nil
}

func (s *VerificationService) getPendingRequest(ctx context.Context, creatorID uuid.UUID) *VerificationRequest {
	return nil
}

func (s *VerificationService) saveVerificationRequest(ctx context.Context, req *VerificationRequest) error {
	return nil
}

func (s *VerificationService) deletePendingRequest(ctx context.Context, creatorID uuid.UUID) error {
	return nil
}
