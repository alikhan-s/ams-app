package flight

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the flight routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	flightGroup := r.Group("/flights")
	{
		flightGroup.POST("/", h.Create)
		flightGroup.GET("/", h.Search)
	}
}
