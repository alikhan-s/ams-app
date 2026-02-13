package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the authentication routes.
func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}
}
