package airportops

import "time"

// Gate represents a terminal gate.
type Gate struct {
	ID         int64  `json:"id"`
	TerminalID int64  `json:"terminal_id"`
	Code       string `json:"code"`
	Status     string `json:"status"`
}

// Baggage represents a baggage item.
type Baggage struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	TagCode   string    `json:"tag_code"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BaggageDetail contains baggage info joined with passenger details.
type BaggageDetail struct {
	ID            int64     `json:"id"`
	TagCode       string    `json:"tag_code"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updated_at"`
	PassengerName string    `json:"passenger_name"` // From users table
	UserID        int64     `json:"user_id"`        // From users table
}

// CreateGateRequest defines the body for creating a gate.
type CreateGateRequest struct {
	TerminalID int64  `json:"terminal_id" binding:"required"`
	Code       string `json:"code" binding:"required"`
	Status     string `json:"status" binding:"required"` // OPEN, CLOSED, MAINTENANCE
}

// CreateBaggageRequest defines the body for checking in baggage.
type CreateBaggageRequest struct {
	TicketID int64 `json:"ticket_id" binding:"required"`
}

// UpdateBaggageRequest defines the body for updating baggage status.
type UpdateBaggageRequest struct {
	Status string `json:"status" binding:"required"`
}
