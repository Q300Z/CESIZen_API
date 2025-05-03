package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddUserRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewUserController(serviceManager)

	rg.GET("/users", controller.GetUsers)
	rg.GET("/users/:id", controller.GetUser)
	rg.GET("/users/search", controller.Search)
	rg.POST("/users", controller.CreateUser)
	rg.PUT("/users/:id", controller.UpdateUser)
	rg.DELETE("/users/:id", controller.DeleteUser)

}
