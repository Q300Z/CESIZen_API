package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddRoleRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewRoleController(serviceManager)

	rg.GET("/roles", controller.GetRoles)
	rg.GET("/roles/:id", controller.GetRole)
	rg.GET("/roles/search", controller.Search)
	rg.POST("/roles", controller.CreateRole)
	rg.PUT("/roles/:id", controller.UpdateRole)
	rg.DELETE("/roles/:id", controller.DeleteRole)

}
