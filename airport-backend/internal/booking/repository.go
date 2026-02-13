package booking

import (
	"airport-system/internal/flight"
	"airport-system/platform/database"
	"context"
	"database/sql"
	"fmt"
)

// Repository handles database interactions for bookings.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new booking repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Create inserts a new ticket.
func (r *Repository) Create(ctx context.Context, ticket *Ticket) (int64, error) {
	var executor database.Executor = r.DB
	if tx := database.GetTx(ctx); tx != nil {
		executor = tx
	}

	query := `
		INSERT INTO tickets (flight_id, passenger_id, seat_no, price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`
	var id int64
	err := executor.QueryRowContext(ctx, query,
		ticket.FlightID,
		ticket.PassengerID,
		ticket.SeatNo,
		ticket.Price,
		ticket.Status,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create ticket: %w", err)
	}
	return id, nil
}

// GetActiveTicketsCount counts active tickets for a flight.
func (r *Repository) GetActiveTicketsCount(ctx context.Context, flightID int64) (int, error) {
	var executor database.Executor = r.DB
	if tx := database.GetTx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT COUNT(*) FROM tickets WHERE flight_id = $1 AND status = 'ACTIVE'`
	var count int
	if err := executor.QueryRowContext(ctx, query, flightID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count active tickets: %w", err)
	}
	return count, nil
}

// GetByPassengerID retrieves all tickets for a specific passenger.
func (r *Repository) GetByPassengerID(ctx context.Context, passengerID int64) ([]Ticket, error) {
	query := `
        SELECT t.id, t.flight_id, t.passenger_id, t.seat_no, t.price, t.status, t.created_at,
               f.id, f.flight_no, f.origin, f.destination, f.departure_time, f.arrival_time, f.status
        FROM tickets t
        JOIN flights f ON t.flight_id = f.id
        WHERE t.passenger_id = $1
    `
	rows, err := r.DB.QueryContext(ctx, query, passengerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var t Ticket
		t.Flight = &flight.Flight{}
		// Scan ticket and embedded flight details
		if err := rows.Scan(
			&t.ID, &t.FlightID, &t.PassengerID, &t.SeatNo, &t.Price, &t.Status, &t.CreatedAt,
			&t.Flight.ID, &t.Flight.FlightNo, &t.Flight.Origin, &t.Flight.Destination,
			&t.Flight.DepartureTime, &t.Flight.ArrivalTime, &t.Flight.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

// GetByID retrieves a ticket by ID.
func (r *Repository) GetByID(ctx context.Context, id int64) (*Ticket, error) {
	query := `SELECT id, flight_id, passenger_id, seat_no, price, status, created_at FROM tickets WHERE id = $1`
	var t Ticket
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.FlightID, &t.PassengerID, &t.SeatNo, &t.Price, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	return &t, nil
}

// Cancel updates ticket status to CANCELLED.
func (r *Repository) Cancel(ctx context.Context, id int64) error {
	query := `UPDATE tickets SET status = 'CANCELLED' WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
