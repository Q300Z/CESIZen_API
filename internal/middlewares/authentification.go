package middlewares

import (
	"cesizen/api/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentification(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		utils.ForbiddenResponse(c, "Invalid token")
		c.Abort()
	}

	c.Set("user", claims)
	c.Next()
}
