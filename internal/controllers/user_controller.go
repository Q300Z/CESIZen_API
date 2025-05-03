package controllers

import (
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.ServiceManager
}

func NewUserController(service *services.ServiceManager) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	// Implement GetUsers logic here
}
func (c *UserController) GetUser(ctx *gin.Context) {
	// Implement GetUser logic here
}
func (c *UserController) Search(ctx *gin.Context) {
	// Implement Search logic here
}
func (c *UserController) CreateUser(ctx *gin.Context) {
	// Implement CreateUser logic here
}
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Implement UpdateUser logic here
}
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Implement DeleteUser logic here
}
