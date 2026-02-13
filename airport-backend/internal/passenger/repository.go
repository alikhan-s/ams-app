package passenger

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository handles database interactions for passengers.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new passenger repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Create inserts a new passenger profile.
func (r *Repository) Create(ctx context.Context, p *Passenger) (int64, error) {
	query := `
		INSERT INTO passengers (user_id, passport_no, phone, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRowContext(ctx, query, p.UserID, p.PassportNo, p.Phone).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create passenger: %w", err)
	}
	return id, nil
}

// GetByUserID retrieves a passenger profile by user ID.
func (r *Repository) GetByUserID(ctx context.Context, userID int64) (*Passenger, error) {
	query := `SELECT id, user_id, passport_no, phone, created_at, updated_at FROM passengers WHERE user_id = $1`
	var p Passenger
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&p.ID, &p.UserID, &p.PassportNo, &p.Phone, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error here, just nil
		}
		return nil, fmt.Errorf("failed to get passenger: %w", err)
	}
	return &p, nil
}
