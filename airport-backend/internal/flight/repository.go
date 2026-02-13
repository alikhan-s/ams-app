package flight

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Repository handles database interactions for flights.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new flight repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Create inserts a new flight into the database.
func (r *Repository) Create(ctx context.Context, f *Flight) (int64, error) {
	query := `
		INSERT INTO flights (flight_no, origin, destination, departure_time, arrival_time, status, version, total_seats, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`

	// Defaults
	if f.Version == 0 {
		f.Version = 1
	}
	if f.TotalSeats == 0 {
		f.TotalSeats = 150
	}

	var id int64
	err := r.DB.QueryRowContext(ctx, query,
		f.FlightNo,
		f.Origin,
		f.Destination,
		f.DepartureTime,
		f.ArrivalTime,
		f.Status,
		f.Version,
		f.TotalSeats,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create flight: %w", err)
	}
	return id, nil
}

// Search retrieves flights based on origin, destination, and date.
func (r *Repository) Search(ctx context.Context, origin, destination string, date time.Time) ([]Flight, error) {
	query := `
		SELECT id, flight_no, origin, destination, gate_id, departure_time, arrival_time, status, version, created_at, updated_at, total_seats
		FROM flights
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if origin != "" {
		query += fmt.Sprintf(" AND origin = $%d", argID)
		args = append(args, origin)
		argID++
	}
	if destination != "" {
		query += fmt.Sprintf(" AND destination = $%d", argID)
		args = append(args, destination)
		argID++
	}

	// Date filtering (Full day range)
	if !date.IsZero() {
		// Start of day
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		// End of day
		endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

		query += fmt.Sprintf(" AND departure_time BETWEEN $%d AND $%d", argID, argID+1)
		args = append(args, startOfDay, endOfDay)
		argID += 2
	}

	query += " ORDER BY departure_time ASC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search flights: %w", err)
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var f Flight
		if err := rows.Scan(
			&f.ID, &f.FlightNo, &f.Origin, &f.Destination, &f.GateID,
			&f.DepartureTime, &f.ArrivalTime, &f.Status, &f.Version,
			&f.CreatedAt, &f.UpdatedAt, &f.TotalSeats,
		); err != nil {
			return nil, fmt.Errorf("failed to scan flight: %w", err)
		}
		flights = append(flights, f)
	}

	return flights, nil
}

// GetByID retrieves a flight by its ID.
func (r *Repository) GetByID(ctx context.Context, id int64) (*Flight, error) {
	query := `
		SELECT id, flight_no, origin, destination, gate_id, departure_time, arrival_time, status, version, created_at, updated_at, total_seats
		FROM flights
		WHERE id = $1
	`
	var f Flight
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.FlightNo, &f.Origin, &f.Destination, &f.GateID,
		&f.DepartureTime, &f.ArrivalTime, &f.Status, &f.Version,
		&f.CreatedAt, &f.UpdatedAt, &f.TotalSeats,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get flight: %w", err)
	}
	return &f, nil
}
