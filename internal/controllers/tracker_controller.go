package controllers

import (
	"cesizen/api/internal/services"

	"github.com/gin-gonic/gin"
)

type TrackerController struct {
	service *services.ServiceManager
}

func NewTrackerController(service *services.ServiceManager) *TrackerController {
	return &TrackerController{service: service}
}

func (c *TrackerController) GetTrackers(ctx *gin.Context) {
	// Implement GetTrackers logic here
}
func (c *TrackerController) GetTracker(ctx *gin.Context) {
	// Implement GetTracker logic here
}
func (c *TrackerController) Search(ctx *gin.Context) {
	// Implement Search logic here
}
func (c *TrackerController) CreateTracker(ctx *gin.Context) {
	// Implement CreateTracker logic here
}
func (c *TrackerController) UpdateTracker(ctx *gin.Context) {
	// Implement UpdateTracker logic here
}
func (c *TrackerController) DeleteTracker(ctx *gin.Context) {
	// Implement DeleteTracker logic here
}
