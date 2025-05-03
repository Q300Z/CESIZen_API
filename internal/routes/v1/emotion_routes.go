package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEmotionRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewEmotionController(serviceManager)

	rg.GET("/emotions", controller.GetEmotions)
	rg.GET("/emotions/:id", controller.GetEmotion)
	rg.GET("/emotions/search", controller.Search)
	rg.POST("/emotions", controller.CreateEmotion)
	rg.PUT("/emotions/:id", controller.UpdateEmotion)
	rg.DELETE("/emotions/:id", controller.DeleteEmotion)

}
