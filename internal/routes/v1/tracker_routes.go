package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddTrackerRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewTrackerController(serviceManager)

	rg.GET("/trackers", controller.GetTrackers)
	rg.GET("/trackers/:id", controller.GetTracker)
	rg.GET("/trackers/search", controller.Search)
	rg.POST("/trackers", controller.CreateTracker)
	rg.PUT("/trackers/:id", controller.UpdateTracker)
	rg.DELETE("/trackers/:id", controller.DeleteTracker)

}
