package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has one of the allowed roles.
// It assumes that the AuthMiddleware has already set the "role" in the context.
func RequireRole(c *gin.Context, allowedRoles ...string) bool {
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
