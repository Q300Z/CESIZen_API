package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/middlewares"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEmotionRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewEmotionController(serviceManager)

	rg.GET("/emotions", middlewares.Authentification, controller.GetEmotions)
	rg.GET("/emotions/:id", middlewares.Authentification, controller.GetEmotion)
	rg.GET("/emotions/search", middlewares.Authentification, controller.Search)
	rg.POST("/emotions", middlewares.Authentification, middlewares.Authorization, controller.CreateEmotion)
	rg.PUT("/emotions/:id", middlewares.Authentification, middlewares.Authorization, controller.UpdateEmotion)
	rg.DELETE("/emotions/:id", middlewares.Authentification, middlewares.Authorization, controller.DeleteEmotion)

	rg.GET("/emotions/base", middlewares.Authentification, controller.GetEmotionBases)
	rg.GET("emotions/base/:id", middlewares.Authentification, middlewares.Authorization, controller.GetEmotionBase)
	rg.PUT("/emotions/base/:id", middlewares.Authentification, middlewares.Authorization, controller.UpdateEmotionBase)
	rg.DELETE("/emotions/base/:id", middlewares.Authentification, middlewares.Authorization, controller.DeleteEmotionBase)
}
