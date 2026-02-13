package flight

import "time"

// Flight represents a flight entity.
type Flight struct {
	ID            int64     `json:"id"`
	FlightNo      string    `json:"flight_no"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	GateID        *int64    `json:"gate_id,omitempty"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Status        string    `json:"status"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TotalSeats    int       `json:"total_seats"`
	BasePrice     float64   `json:"base_price"`
}

// CreateFlightParams defines the parameters for creating a new flight.
type CreateFlightParams struct {
	FlightNo      string `json:"flight_no" binding:"required"`
	Origin        string `json:"origin" binding:"required,len=3"`
	Destination   string `json:"destination" binding:"required,len=3"`
	DepartureTime string `json:"departure_time" binding:"required"` // Format: RFC3339
	ArrivalTime   string `json:"arrival_time" binding:"required"`   // Format: RFC3339
}

// SearchParams defines criteria for searching flights.
type SearchParams struct {
	Origin      string `form:"origin"`
	Destination string `form:"destination"`
	Date        string `form:"date"` // Format: YYYY-MM-DD
}
