package routes

import (
	v1 "cesizen/api/internal/routes/v1"
	"cesizen/api/internal/services"

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
}
