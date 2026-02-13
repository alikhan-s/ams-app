package airportops

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the airport ops routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	opsGroup := r.Group("/ops")
	opsGroup.Use(authMiddleware)
	{
		opsGroup.POST("/gates", h.CreateGate)
		opsGroup.GET("/gates", h.ListGates)
		opsGroup.POST("/baggage", h.CheckInBaggage)
		opsGroup.GET("/baggage", h.ListBaggage)
		opsGroup.PATCH("/baggage/:id", h.UpdateBaggage)
	}
}
