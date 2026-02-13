package airportops

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository handles database interactions for airport operations.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new airport ops repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// CreateGate inserts a new gate.
func (r *Repository) CreateGate(ctx context.Context, gate *Gate) (int64, error) {
	query := `
		INSERT INTO gates (terminal_id, code, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRowContext(ctx, query, gate.TerminalID, gate.Code, gate.Status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create gate: %w", err)
	}
	return id, nil
}

// ListGates returns all gates.
func (r *Repository) ListGates(ctx context.Context) ([]Gate, error) {
	query := `SELECT id, terminal_id, code, status FROM gates`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list gates: %w", err)
	}
	defer rows.Close()

	var gates []Gate
	for rows.Next() {
		var g Gate
		if err := rows.Scan(&g.ID, &g.TerminalID, &g.Code, &g.Status); err != nil {
			return nil, fmt.Errorf("failed to scan gate: %w", err)
		}
		gates = append(gates, g)
	}
	return gates, nil
}

// CreateBaggage inserts a new baggage item.
func (r *Repository) CreateBaggage(ctx context.Context, bag *Baggage) (int64, error) {
	query := `
		INSERT INTO baggage (ticket_id, tag_code, status, updated_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRowContext(ctx, query, bag.TicketID, bag.TagCode, bag.Status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create baggage: %w", err)
	}
	return id, nil
}

// UpdateBaggageStatus updates the status of a baggage item.
func (r *Repository) UpdateBaggageStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE baggage SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update baggage status: %w", err)
	}
	return nil
}

// GetBaggageByID retrieves baggage by ID (helper for service).
func (r *Repository) GetBaggageByID(ctx context.Context, id int64) (*Baggage, error) {
	query := `SELECT id, ticket_id, tag_code, status, updated_at FROM baggage WHERE id = $1`
	var b Baggage
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.TicketID, &b.TagCode, &b.Status, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get baggage: %w", err)
	}
	return &b, nil
}

// ListAllWithPassengerInfo retrieves all baggage with passenger details (For Admin/Staff).
func (r *Repository) ListAllWithPassengerInfo(ctx context.Context) ([]BaggageDetail, error) {
	query := `
		SELECT b.id, b.tag_code, b.status, b.updated_at, u.full_name, u.id
		FROM baggage b
		JOIN tickets t ON b.ticket_id = t.id
		JOIN passengers p ON t.passenger_id = p.id
		JOIN users u ON p.user_id = u.id
		ORDER BY b.updated_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all baggage: %w", err)
	}
	defer rows.Close()

	var details []BaggageDetail
	for rows.Next() {
		var b BaggageDetail
		if err := rows.Scan(&b.ID, &b.TagCode, &b.Status, &b.UpdatedAt, &b.PassengerName, &b.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan baggage detail: %w", err)
		}
		details = append(details, b)
	}
	return details, nil
}

// GetByTicketID retrieves baggage items by ticket ID (For Booking module).
func (r *Repository) GetByTicketID(ctx context.Context, ticketID int64) ([]Baggage, error) {
	query := `SELECT id, ticket_id, tag_code, status, updated_at FROM baggage WHERE ticket_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get baggage by ticket: %w", err)
	}
	defer rows.Close()

	var bags []Baggage
	for rows.Next() {
		var b Baggage
		if err := rows.Scan(&b.ID, &b.TicketID, &b.TagCode, &b.Status, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan baggage: %w", err)
		}
		bags = append(bags, b)
	}
	return bags, nil
}
