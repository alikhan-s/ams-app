package booking

import (
	"airport-system/internal/airportops"
	"airport-system/internal/flight"
	"airport-system/internal/passenger"
	"airport-system/platform/database"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Service handles booking business logic.
type Service struct {
	repo        *Repository
	flightRepo  *flight.Repository
	txManager   database.TxManager
	opsService  *airportops.Service
	passService *passenger.Service
	log         *slog.Logger
}

// NewService creates a new booking service.
func NewService(repo *Repository, flightRepo *flight.Repository, txManager database.TxManager, opsService *airportops.Service, passService *passenger.Service, log *slog.Logger) *Service {
	return &Service{
		repo:        repo,
		flightRepo:  flightRepo,
		txManager:   txManager,
		opsService:  opsService,
		passService: passService,
		log:         log,
	}
}

// BookTicket books a ticket for a user on a flight transactionally.
func (s *Service) BookTicket(ctx context.Context, userID int64, req BookingRequest) (*Ticket, error) {
	var ticket *Ticket

	// 0. Resolve Passenger ID
	// Check if passenger profile exists
	passProfile, err := s.passService.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check passenger profile: %w", err)
	}

	var passengerID int64
	if passProfile != nil {
		passengerID = passProfile.ID
	} else {
		// Profile does not exist, require passport info
		if req.PassportNo == "" {
			return nil, errors.New("passenger profile required: please provide passport_no and phone")
		}
		// Create new profile
		newProfile, err := s.passService.CreateProfile(ctx, userID, req.PassportNo, req.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to create passenger profile: %w", err)
		}
		passengerID = newProfile.ID
	}

	err = s.txManager.Run(ctx, func(ctx context.Context) error {
		// 1. Get Flight details to check capacity
		f, err := s.flightRepo.GetByID(ctx, req.FlightID)
		if err != nil {
			return fmt.Errorf("failed to get flight: %w", err)
		}
		if f == nil {
			return errors.New("flight not found")
		}

		// 2. Check capacity (Optimistic locking via count)
		activeCount, err := s.repo.GetActiveTicketsCount(ctx, req.FlightID)
		if err != nil {
			return fmt.Errorf("failed to count tickets: %w", err)
		}

		if activeCount >= f.TotalSeats {
			return errors.New("flight is full")
		}

		// 3. Create Ticket using resolved PassengerID
		ticket = &Ticket{
			FlightID:    req.FlightID,
			PassengerID: passengerID, // Use PassengerID, not UserID
			Price:       f.BasePrice,
			Status:      "ACTIVE",
		}

		id, err := s.repo.Create(ctx, ticket)
		if err != nil {
			return err
		}
		ticket.ID = id
		ticket.Flight = f // Attach flight details for response
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// GetMyBookings returns bookings for the current user.
func (s *Service) GetMyBookings(ctx context.Context, userID int64) ([]Ticket, error) {
	// First get passenger ID for user
	passProfile, err := s.passService.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get passenger profile: %w", err)
	}
	if passProfile == nil {
		return []Ticket{}, nil // No profile means no bookings
	}

	// Fetch bookings using passenger ID
	return s.repo.GetByPassengerID(ctx, passProfile.ID)
}

// CancelTicket cancels a user's ticket.
func (s *Service) CancelTicket(ctx context.Context, userID, ticketID int64) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}

	// Verify ownership via passenger profile
	passProfile, err := s.passService.GetProfile(ctx, userID)
	if err != nil {
		return err
	}
	if passProfile == nil || ticket.PassengerID != passProfile.ID {
		return errors.New("unauthorized to cancel this ticket")
	}

	if ticket.Status == "CANCELLED" {
		return errors.New("ticket is already cancelled")
	}

	return s.repo.Cancel(ctx, ticketID)
}

// GetUserBaggage returns all baggage for the current user across all bookings.
func (s *Service) GetUserBaggage(ctx context.Context, userID, ticketID int64) ([]airportops.Baggage, error) {
	// 1. Get Passenger Profile
	passProfile, err := s.passService.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get passenger profile: %w", err)
	}
	if passProfile == nil {
		return []airportops.Baggage{}, nil
	}

	// 2. Get Bookings
	// If ticketID is provided, verify it belongs to the user
	var targetBookingIDs []int64

	if ticketID != 0 {
		ticket, err := s.repo.GetByID(ctx, ticketID)
		if err != nil {
			return nil, err
		}
		if ticket == nil {
			return nil, errors.New("ticket not found")
		}
		if ticket.PassengerID != passProfile.ID {
			return nil, errors.New("unauthorized to view baggage for this ticket")
		}
		targetBookingIDs = append(targetBookingIDs, ticketID)
	} else {
		// Fetch all bookings
		bookings, err := s.repo.GetByPassengerID(ctx, passProfile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get bookings: %w", err)
		}
		for _, b := range bookings {
			targetBookingIDs = append(targetBookingIDs, b.ID)
		}
	}

	// 3. Aggregate Baggage
	var allBaggage []airportops.Baggage
	for _, bookingID := range targetBookingIDs {
		bags, err := s.opsService.GetBaggageByTicketID(ctx, bookingID)
		if err != nil {
			return nil, fmt.Errorf("failed to get baggage for ticket %d: %w", bookingID, err)
		}
		allBaggage = append(allBaggage, bags...)
	}

	return allBaggage, nil
}
