package routes

import (
	v1 "cesizen/api/internal/routes/v1"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetRoutes(r *gin.Engine, serviceManager *services.ServiceManager) {
	v1Group := r.Group("/v1")
	v1Group.Static("/uploads", "./uploads")
	v1.AddArticleRoutes(v1Group, serviceManager)
	v1.AddAuthRoutes(v1Group, serviceManager)
	v1.AddEmotionRoutes(v1Group, serviceManager)
	v1.AddTrackerRoutes(v1Group, serviceManager)
	v1.AddUserRoutes(v1Group, serviceManager)
	v1Group.GET("/version", func(ctx *gin.Context) {
		// Lire le fichier de version
		// et renvoyer la version de l'API

		ctx.JSON(200, gin.H{
			"version": utils.GetVersion(),
			"mode":    utils.GetEnv("GIN_MODE", "ERROR"),
			"message": "API is running",
		})
	})
}
