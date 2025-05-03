package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/middlewares"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddTrackerRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewTrackerController(serviceManager)

	rg.GET("/trackers", middlewares.Authentification, controller.GetTrackers)
	rg.GET("/trackers/:id", middlewares.Authentification, controller.GetTracker)
	rg.GET("/trackers/search", middlewares.Authentification, controller.Search)
	rg.POST("/trackers", middlewares.Authentification, controller.CreateTracker)
	rg.PUT("/trackers/:id", middlewares.Authentification, controller.UpdateTracker)
	rg.DELETE("/trackers/:id", middlewares.Authentification, controller.DeleteTracker)

}
