package booking

import (
	"airport-system/internal/flight"
	"time"
)

// Ticket represents a booked flight ticket.
type Ticket struct {
	ID          int64          `json:"id"`
	FlightID    int64          `json:"flight_id"`
	Flight      *flight.Flight `json:"flight,omitempty"` // For joining flight details
	PassengerID int64          `json:"passenger_id"`
	SeatNo      *string        `json:"seat_no"`
	Price       float64        `json:"price"`
	Status      string         `json:"status"` // ACTIVE, CANCELLED
	CreatedAt   time.Time      `json:"created_at"`
}

// BookingRequest defines the body for booking a ticket.
type BookingRequest struct {
	FlightID   int64  `json:"flight_id" binding:"required"`
	PassportNo string `json:"passport_no"` // Optional: required only if profile doesn't exist
	Phone      string `json:"phone"`       // Optional
}
