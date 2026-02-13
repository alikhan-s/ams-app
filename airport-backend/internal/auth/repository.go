package auth

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository handles database operations for users.
type Repository struct {
	DB *sql.DB
}

// NewRepository creates a new user repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// CreateUser inserts a new user into the database and returns the ID.
func (r *Repository) CreateUser(ctx context.Context, user *User) (int64, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRowContext(ctx, query, user.FullName, user.Email, user.PasswordHash, user.Role).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	user := &User{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}
