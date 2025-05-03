package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TrackerController struct {
	service *services.ServiceManager
}

func NewTrackerController(service *services.ServiceManager) *TrackerController {
	return &TrackerController{service: service}
}

// GET /trackers
func (c *TrackerController) GetTrackers(ctx *gin.Context) {
	trackers, err := c.service.Client.Tracker.FindMany().With(
		db.Tracker.Emotion.Fetch().With(db.Emotion.EmotionBase.Fetch()),
	).Exec(c.service.Ctx)

	if err != nil {
		log.Println("Error fetching trackers:", err)
		utils.ErrorResponse(ctx, 500, "Failed to fetch trackers", err)
		return
	}
	utils.SuccessResponse(ctx, "Trackers fetched successfully", trackers)
}

// GET /tracker/:id
func (c *TrackerController) GetTracker(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid tracker ID", err)
		return
	}

	tracker, err := c.service.Client.Tracker.FindUnique(
		db.Tracker.ID.Equals(id),
	).With(
		db.Tracker.Emotion.Fetch().With(db.Emotion.EmotionBase.Fetch()),
	).Exec(c.service.Ctx)

	if err != nil || tracker == nil {
		utils.ErrorResponse(ctx, 404, "Tracker not found", err)
		return
	}

	utils.SuccessResponse(ctx, "Tracker fetched successfully", tracker)
}

// GET /trackers/search?q=keyword
func (c *TrackerController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if strings.TrimSpace(query) == "" {
		utils.ErrorResponse(ctx, 400, "Search query is required", nil)
		return
	}

	trackers, err := c.service.Client.Tracker.FindMany(
		db.Tracker.Description.Contains(query),
	).With(
		db.Tracker.Emotion.Fetch().With(db.Emotion.EmotionBase.Fetch()),
	).Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Search failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Search results", trackers)
}

// POST /tracker
func (c *TrackerController) CreateTracker(ctx *gin.Context) {
	var input struct {
		Description string `json:"description"`
		UserID      int    `json:"userId" binding:"required"`
		EmotionID   int    `json:"emotionId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	create := c.service.Client.Tracker.CreateOne(
		db.Tracker.User.Link(
			db.User.ID.Equals(input.UserID),
		),
		db.Tracker.Emotion.Link(
			db.Emotion.ID.Equals(input.EmotionID),
		),
	)

	tracker, err := create.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to create tracker", err)
		return
	}

	utils.SuccessResponse(ctx, "Tracker created", tracker)
}

// PUT /trackers/:id
func (c *TrackerController) UpdateTracker(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid tracker ID", err)
		return
	}

	var input struct {
		Description *string `json:"description"`
		EmotionID   *int    `json:"emotionId"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	update := c.service.Client.Tracker.FindUnique(
		db.Tracker.ID.Equals(id),
	).Update(
		db.Tracker.Description.SetIfPresent(input.Description),
		db.Tracker.Emotion.Link(
			db.Emotion.ID.Equals(*input.EmotionID),
		),
	)

	tracker, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update tracker", err)
		return
	}

	utils.SuccessResponse(ctx, "Tracker updated", tracker)
}

// DELETE /trackers/:id
func (c *TrackerController) DeleteTracker(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid tracker ID", err)
		return
	}

	_, err = c.service.Client.Tracker.FindUnique(
		db.Tracker.ID.Equals(id),
	).Delete().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to delete tracker", err)
		return
	}

	utils.SuccessResponse(ctx, "Tracker deleted successfully", nil)
}
