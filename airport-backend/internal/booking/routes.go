package booking

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the booking routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	bookingGroup := r.Group("/bookings")
	bookingGroup.Use(authMiddleware)
	{
		bookingGroup.POST("/", h.Book)
		bookingGroup.GET("/my", h.GetMy)
		bookingGroup.POST("/:id/cancel", h.Cancel)
		bookingGroup.GET("/baggage", h.GetMyBaggage)
	}
}
