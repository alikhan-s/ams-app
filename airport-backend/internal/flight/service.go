package flight

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// Service handles business logic for flights.
type Service struct {
	repo *Repository
	log  *slog.Logger
}

// NewService creates a new flight service.
func NewService(repo *Repository, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// CreateFlight validates and creates a new flight.
func (s *Service) CreateFlight(ctx context.Context, params CreateFlightParams) (*Flight, error) {
	// Parse times
	depTime, err := time.Parse(time.RFC3339, params.DepartureTime)
	if err != nil {
		return nil, errors.New("invalid departure_time format (expected RFC3339)")
	}
	arrTime, err := time.Parse(time.RFC3339, params.ArrivalTime)
	if err != nil {
		return nil, errors.New("invalid arrival_time format (expected RFC3339)")
	}

	// Validation
	if !arrTime.After(depTime) {
		return nil, errors.New("arrival_time must be after departure_time")
	}
	if params.Origin == params.Destination {
		return nil, errors.New("origin and destination cannot be the same")
	}

	flight := &Flight{
		FlightNo:      params.FlightNo,
		Origin:        params.Origin,
		Destination:   params.Destination,
		DepartureTime: depTime,
		ArrivalTime:   arrTime,
		Status:        "SCHEDULED",
	}

	id, err := s.repo.Create(ctx, flight)
	if err != nil {
		return nil, err
	}
	flight.ID = id
	return flight, nil
}

// SearchFlights searches for flights based on criteria.
func (s *Service) SearchFlights(ctx context.Context, origin, destination, dateStr string) ([]Flight, error) {
	var searchDate time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, errors.New("invalid date format (expected YYYY-MM-DD)")
		}
		searchDate = parsedDate
	}

	return s.repo.Search(ctx, origin, destination, searchDate)
}

// GetByID retrieves a flight by its ID.
func (s *Service) GetByID(ctx context.Context, id int64) (*Flight, error) {
	return s.repo.GetByID(ctx, id)
}
