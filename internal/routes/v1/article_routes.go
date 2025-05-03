package v1

import (
	"cesizen/api/internal/controllers"
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

func AddArticleRoutes(rg *gin.RouterGroup, serviceManager *services.ServiceManager) {

	controller := controllers.NewArticleController(serviceManager)

	rg.GET("/articles", controller.GetArticles)
	rg.GET("/articles/:id", controller.GetArticle)
	rg.GET("/articles/search", controller.Search)
	rg.POST("/articles", controller.CreateArticle)
	rg.PUT("/articles/:id", controller.UpdateArticle)
	rg.DELETE("/articles/:id", controller.DeleteArticle)

}
