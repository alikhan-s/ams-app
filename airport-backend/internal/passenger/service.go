package passenger

import (
	"context"
	"log/slog"
)

// Service handles business logic for passengers.
type Service struct {
	repo *Repository
	log  *slog.Logger
}

// NewService creates a new passenger service.
func NewService(repo *Repository, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// GetProfile retrieves a passenger profile by user ID.
func (s *Service) GetProfile(ctx context.Context, userID int64) (*Passenger, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// CreateProfile creates a new passenger profile.
func (s *Service) CreateProfile(ctx context.Context, userID int64, passport, phone string) (*Passenger, error) {
	p := &Passenger{
		UserID:     userID,
		PassportNo: passport,
		Phone:      phone,
	}

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}
