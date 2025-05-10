package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/helpers"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type EmotionController struct {
	service *services.ServiceManager
}

func NewEmotionController(service *services.ServiceManager) *EmotionController {
	return &EmotionController{service: service}
}

// GET /emotions
func (c *EmotionController) GetEmotions(ctx *gin.Context) {
	emotions, err := c.service.Client.Emotion.FindMany().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to fetch emotions", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotions fetched", emotions)
}

// GET /emotions/:id
func (c *EmotionController) GetEmotion(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion ID", err)
		return
	}

	emotion, err := c.service.Client.Emotion.FindUnique(
		db.Emotion.ID.Equals(id),
	).With(db.Emotion.EmotionBase.Fetch()).Exec(c.service.Ctx)

	if err != nil || emotion == nil {
		utils.ErrorResponse(ctx, 404, "Emotion not found", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion fetched", emotion)
}

// GET /emotions/search?q=mot
func (c *EmotionController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if strings.TrimSpace(query) == "" {
		utils.ErrorResponse(ctx, 400, "Query string required", nil)
		return
	}

	results, err := c.service.Client.Emotion.FindMany(
		db.Emotion.Name.Contains(query),
	).With(db.Emotion.EmotionBase.Fetch()).Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Search failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Search results", results)
}

// POST /emotions
func (c *EmotionController) CreateEmotion(ctx *gin.Context) {
	// Parse form-data (limite 10 Mo)
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(ctx, 400, "Failed to parse form", err)
		return
	}

	name := ctx.PostForm("name")
	emotionBaseIDStr := ctx.PostForm("emotionBaseID")

	if name == "" {
		utils.ErrorResponse(ctx, 400, "Name is required", nil)
		return
	}

	// Gérer l'image via le helper
	var imageURL *string
	fileHeader, err := ctx.FormFile("image")
	if err == nil {
		_, url, err := helpers.SaveImage(fileHeader)
		if err != nil {
			utils.ErrorResponse(ctx, 500, "Image upload failed", err)
			return
		}
		imageURL = &url
	}

	// S'il n'y a pas d'EmotionBaseID, on crée une EmotionBase
	if emotionBaseIDStr == "" {
		baseCreate := c.service.Client.EmotionBase.CreateOne(
			db.EmotionBase.Name.Set(name),
		)
		base, err := baseCreate.Exec(c.service.Ctx)
		if err != nil {
			utils.ErrorResponse(ctx, 500, "Failed to create emotion base", err)
			return
		}
		utils.SuccessResponse(ctx, "Emotion base created", base)
		return
	}

	// Convertir EmotionBaseID
	emotionBaseID, err := strconv.Atoi(emotionBaseIDStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotionBaseId", err)
		return
	}

	// Créer une Emotion liée à une EmotionBase
	emotionCreate := c.service.Client.Emotion.CreateOne(
		db.Emotion.Name.Set(name),
		db.Emotion.Image.SetIfPresent(imageURL),
		db.Emotion.EmotionBase.Link(
			db.EmotionBase.ID.Equals(emotionBaseID),
		),
	)

	emotion, err := emotionCreate.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to create emotion", err)
		return
	}

	utils.CreatedResponse(ctx, "Emotion created", emotion)
}

// PUT /emotions/:id
func (c *EmotionController) UpdateEmotion(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion ID", err)
		return
	}

	// Parse form-data (limite 10 Mo)
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(ctx, 400, "Failed to parse form", err)
		return
	}

	// Récupérer les champs du formulaire
	name := ctx.PostForm("name")
	emotionBaseIDStr := ctx.PostForm("emotionBaseId")

	// Gérer l'image via le helper si elle est présente
	var imageURL *string
	fileHeader, err := ctx.FormFile("image")
	if err == nil {
		_, url, err := helpers.SaveImage(fileHeader)
		if err != nil {
			utils.ErrorResponse(ctx, 500, "Image upload failed", err)
			return
		}
		imageURL = &url
	}
	update := c.service.Client.Emotion.FindUnique(
		db.Emotion.ID.Equals(id),
	).Update(
		db.Emotion.Name.Set(name),
		db.Emotion.Image.SetIfPresent(imageURL),
	)

	if emotionBaseIDStr != "" {
		emotionBaseID, err := strconv.Atoi(emotionBaseIDStr)
		if err != nil {
			utils.ErrorResponse(ctx, 400, "Invalid emotionBaseId", err)
			return
		}
		update = c.service.Client.Emotion.FindUnique(
			db.Emotion.ID.Equals(id),
		).Update(
			db.Emotion.Name.Set(name),
			db.Emotion.Image.SetIfPresent(imageURL),
			db.Emotion.EmotionBase.Link(
				db.EmotionBase.ID.Equals(emotionBaseID),
			),
		)
	}

	emotion, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update emotion", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion updated", emotion)
}

// DELETE /emotions/:id
func (c *EmotionController) DeleteEmotion(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion ID", err)
		return
	}

	_, err = c.service.Client.Emotion.FindUnique(
		db.Emotion.ID.Equals(id),
	).Delete().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to delete emotion", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion deleted", nil)
}

// GET /emotions/base
func (c *EmotionController) GetEmotionBases(ctx *gin.Context) {
	bases, err := c.service.Client.EmotionBase.FindMany().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to fetch emotion bases", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion bases fetched", bases)
}

// GET /emotions/base/:id
func (c *EmotionController) GetEmotionBase(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion base ID", err)
		return
	}

	emotion, err := c.service.Client.EmotionBase.FindUnique(
		db.EmotionBase.ID.Equals(id),
	).With(db.EmotionBase.Emotions.Fetch()).Exec(c.service.Ctx)

	if err != nil || emotion == nil {
		utils.ErrorResponse(ctx, 404, "Emotion base not found", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion base fetched", emotion)
}

// PUT /emotions/base/:id
func (c *EmotionController) UpdateEmotionBase(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion base ID", err)
		return
	}

	var input struct {
		Name *string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	update := c.service.Client.EmotionBase.FindUnique(
		db.EmotionBase.ID.Equals(id),
	).Update(
		db.EmotionBase.Name.SetIfPresent(input.Name),
	)

	emotion, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update emotion base", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion base updated", emotion)
}

// DELETE /emotions/base/:id
func (c *EmotionController) DeleteEmotionBase(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid emotion base ID", err)
		return
	}

	_, err = c.service.Client.EmotionBase.FindUnique(
		db.EmotionBase.ID.Equals(id),
	).Delete().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to delete emotion", err)
		return
	}

	utils.SuccessResponse(ctx, "Emotion base deleted", nil)
}
