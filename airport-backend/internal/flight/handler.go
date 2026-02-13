package flight

import (
	"airport-system/internal/auth"
	"net/http"
	"strconv"

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
	if !auth.RequireRole(c, "ADMIN") {
		return
	}

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

// GetByID handles getting a flight by ID.
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	flight, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if flight == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}
