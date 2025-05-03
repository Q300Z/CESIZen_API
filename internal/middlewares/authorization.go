package middlewares

import (
	"cesizen/api/internal/models"
	"cesizen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		utils.ForbiddenResponse(c, "Forbidden")
		c.Abort()
		return
	}

	// models.JWTClaims
	userClaims := user.(models.JWTClaims)

	// Check if the user has the required role
	if userRole := userClaims.Role; userRole != "admin" {
		utils.ForbiddenResponse(c, "Forbidden")
		c.Abort()
		return
	}

	c.Next()
}
