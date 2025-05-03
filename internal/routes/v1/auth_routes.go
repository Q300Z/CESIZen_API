package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/middlewares"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddAuthRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewAuthController(serviceManager)

	rg.POST("/login", controller.Login)
	rg.POST("/register", controller.Register)
	rg.GET("/logout", middlewares.Authentification, controller.Logout)
	rg.GET("/logout-admin", middlewares.Authentification, middlewares.Authorization, controller.Logout)
	rg.POST("/reset-password", controller.ResetPassword)
	rg.POST("/change-password", middlewares.Authentification, controller.ChangePassword)
}
