package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/middlewares"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddUserRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewUserController(serviceManager)

	rg.GET("/users", middlewares.Authentification, middlewares.Authorization, controller.GetUsers)
	rg.GET("/users/:id", middlewares.Authentification, middlewares.Authorization, controller.GetUser)
	rg.GET("/users/search", middlewares.Authentification, middlewares.Authorization, controller.Search)
	rg.POST("/users", middlewares.Authentification, middlewares.Authorization, controller.CreateUser)
	rg.PUT("/users/:id", middlewares.Authentification, controller.UpdateUser)
	rg.DELETE("/users/:id", middlewares.Authentification, middlewares.Authorization, controller.DeleteUser)
	rg.DELETE("/users/me", middlewares.Authentification, controller.DeleteCurrentUser)

}
