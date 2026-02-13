package airportops

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
)

// Service handles business logic for airport operations.
type Service struct {
	repo *Repository
	log  *slog.Logger
}

// NewService creates a new airport ops service.
func NewService(repo *Repository, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// CreateGate creates a new gate.
func (s *Service) CreateGate(ctx context.Context, req CreateGateRequest) (*Gate, error) {
	gate := &Gate{
		TerminalID: req.TerminalID,
		Code:       req.Code,
		Status:     req.Status,
	}

	id, err := s.repo.CreateGate(ctx, gate)
	if err != nil {
		return nil, err
	}
	gate.ID = id
	return gate, nil
}

// ListGates lists all gates.
func (s *Service) ListGates(ctx context.Context) ([]Gate, error) {
	return s.repo.ListGates(ctx)
}

// CheckInBaggage generates a tag and checks in baggage.
func (s *Service) CheckInBaggage(ctx context.Context, ticketID int64) (*Baggage, error) {
	// Generate unique tag
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate tag: %w", err)
	}
	tagCode := fmt.Sprintf("BAG-%s", hex.EncodeToString(bytes))

	bag := &Baggage{
		TicketID: ticketID,
		TagCode:  tagCode,
		Status:   "RECEIVED",
	}

	id, err := s.repo.CreateBaggage(ctx, bag)
	if err != nil {
		return nil, err
	}
	bag.ID = id
	return bag, nil
}

// UpdateBaggage updates the status of a baggage item.
func (s *Service) UpdateBaggage(ctx context.Context, id int64, status string) (*Baggage, error) {
	if err := s.repo.UpdateBaggageStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.repo.GetBaggageByID(ctx, id)
}

// ListAllBaggage returns detailed baggage info for Staff/Admin.
func (s *Service) ListAllBaggage(ctx context.Context) ([]BaggageDetail, error) {
	return s.repo.ListAllWithPassengerInfo(ctx)
}

// GetBaggageByTicketID returns baggage for a specific ticket (Used by Booking module).
func (s *Service) GetBaggageByTicketID(ctx context.Context, ticketID int64) ([]Baggage, error) {
	return s.repo.GetByTicketID(ctx, ticketID)
}
