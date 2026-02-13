package flight

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler manages HTTP requests for flights.
type Handler struct {
	Service *Service
}

// NewHandler creates a new flight handler.
func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

// Create handles flight creation.
func (h *Handler) Create(c *gin.Context) {
	var req CreateFlightParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flight, err := h.Service.CreateFlight(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flight)
}

// Search handles flight search.
func (h *Handler) Search(c *gin.Context) {
	var params SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flights, err := h.Service.SearchFlights(c.Request.Context(), params.Origin, params.Destination, params.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flights)
}
