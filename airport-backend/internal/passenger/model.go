package passenger

import "time"

// Passenger represents a passenger profile.
type Passenger struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	PassportNo string    `json:"passport_no"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
