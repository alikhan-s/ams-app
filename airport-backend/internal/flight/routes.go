package flight

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the flight routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	flightGroup := r.Group("/flights")
	{
		flightGroup.GET("", h.Search) 
		
		flightGroup.GET("/:id", h.GetByID)

		// Protected routes
		flightGroup.POST("", authMiddleware, h.Create)
	}
}