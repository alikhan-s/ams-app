package airportops

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler manages HTTP requests for airport operations.
type Handler struct {
	Service *Service
}

// NewHandler creates a new airport ops handler.
func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

// requireRole checks if the user has one of the allowed roles.
func (h *Handler) requireRole(c *gin.Context, allowedRoles ...string) bool {
	userRole := c.GetString("role")
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	c.Abort()
	return false
}

// CreateGate handles gate creation (ADMIN only).
func (h *Handler) CreateGate(c *gin.Context) {
	if !h.requireRole(c, "ADMIN") {
		return
	}

	var req CreateGateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gate, err := h.Service.CreateGate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gate)
}

// ListGates lists all gates (STAFF, ADMIN).
func (h *Handler) ListGates(c *gin.Context) {
	if !h.requireRole(c, "STAFF", "ADMIN") {
		return
	}

	gates, err := h.Service.ListGates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gates)
}

// CheckInBaggage handles baggage check-in (STAFF, ADMIN).
func (h *Handler) CheckInBaggage(c *gin.Context) {
	if !h.requireRole(c, "STAFF", "ADMIN") {
		return
	}

	var req CreateBaggageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bag, err := h.Service.CheckInBaggage(c.Request.Context(), req.TicketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bag)
}

// UpdateBaggage handles baggage status update (STAFF, ADMIN).
func (h *Handler) UpdateBaggage(c *gin.Context) {
	if !h.requireRole(c, "STAFF", "ADMIN") {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateBaggageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bag, err := h.Service.UpdateBaggage(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bag)
}

// ListBaggage lists all baggage with passenger info (STAFF, ADMIN).
func (h *Handler) ListBaggage(c *gin.Context) {
	if !h.requireRole(c, "STAFF", "ADMIN") {
		return
	}

	baggageList, err := h.Service.ListAllBaggage(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, baggageList)
}
