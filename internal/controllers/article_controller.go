package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/models"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type ArticleController struct {
	service *services.ServiceManager
}

func NewArticleController(service *services.ServiceManager) *ArticleController {
	return &ArticleController{service: service}
}

// GET /articles
func (c *ArticleController) GetArticles(ctx *gin.Context) {
	articles, err := c.service.Client.Article.
		FindMany().With(db.Article.User.Fetch()).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to fetch articles", err)
		return
	}

	var response []models.ArticleResponse
	for _, article := range articles {
		desc, ok := article.Description()
		if ok {
			desc = "(pas de description)"
		}

		response = append(response, models.ArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			Description: desc,
			Content:     article.Content,
			CreatedAt:   article.CreateAt,
			UpdatedAt:   article.UpdateAt,
			User: &models.UserResponse{
				ID:        article.User().ID,
				Name:      article.User().Name,
				Email:     article.User().Email,
				Role:      string(article.User().Role),
				CreatedAt: article.User().CreateAt,
				UpdatedAt: article.User().UpdateAt,
			},
		})
	}

	utils.SuccessResponse(ctx, "Articles fetched successfully", response)
}

// GET /article/:id
func (c *ArticleController) GetArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid article ID", err)
		return
	}

	article, err := c.service.Client.Article.FindUnique(
		db.Article.ID.Equals(id),
	).With(db.Article.User.Fetch()).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "Article not found", err)
		return
	}

	var response models.ArticleResponse
	desc, ok := article.Description()
	if ok {
		desc = "(pas de description)"
	}
	response = models.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Description: desc,
		Content:     article.Content,
		CreatedAt:   article.CreateAt,
		UpdatedAt:   article.UpdateAt,
		User: &models.UserResponse{
			ID:        article.User().ID,
			Name:      article.User().Name,
			Email:     article.User().Email,
			Role:      string(article.User().Role),
			CreatedAt: article.User().CreateAt,
			UpdatedAt: article.User().UpdateAt,
		},
	}

	utils.SuccessResponse(ctx, "Article fetched successfully", response)
}

// GET /articles/search?q=keyword
func (c *ArticleController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if strings.TrimSpace(query) == "" {
		utils.ErrorResponse(ctx, 400, "Search query is required", nil)
		return
	}

	articles, err := c.service.Client.Article.FindMany(
		db.Article.Title.Contains(query),
	).With(db.Article.User.Fetch()).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Search failed", err)
		return
	}

	var response []models.ArticleResponse
	for _, article := range articles {
		desc, ok := article.Description()
		if ok {
			desc = "(pas de description)"
		}

		response = append(response, models.ArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			Description: desc,
			Content:     article.Content,
			CreatedAt:   article.CreateAt,
			UpdatedAt:   article.UpdateAt,
			User: &models.UserResponse{
				ID:        article.User().ID,
				Name:      article.User().Name,
				Email:     article.User().Email,
				Role:      string(article.User().Role),
				CreatedAt: article.User().CreateAt,
				UpdatedAt: article.User().UpdateAt,
			},
		})
	}

	utils.SuccessResponse(ctx, "Search results", response)
}

// POST /article
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Content     string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	// Get the user from the context
	user, exists := ctx.Get("user")
	if !exists {
		utils.ForbiddenResponse(ctx, "Forbidden")
		ctx.Abort()
		return
	}

	// models.JWTClaims
	userClaims := user.(models.JWTClaims)

	article, err := c.service.Client.Article.CreateOne(
		db.Article.Title.Set(input.Title),
		db.Article.Content.Set(input.Content),
		db.Article.User.Link(
			db.User.ID.Equals(int(userClaims.UserID)),
		),
		db.Article.Description.Set(input.Description),
	).With(db.Article.User.Fetch()).Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to create article", err)
		return
	}

	var response models.ArticleResponse
	desc, ok := article.Description()
	if ok {
		desc = "(pas de description)"
	}
	response = models.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Description: desc,
		Content:     article.Content,
		CreatedAt:   article.CreateAt,
		UpdatedAt:   article.UpdateAt,
		User: &models.UserResponse{
			ID:        article.User().ID,
			Name:      article.User().Name,
			Email:     article.User().Email,
			Role:      string(article.User().Role),
			CreatedAt: article.User().CreateAt,
			UpdatedAt: article.User().UpdateAt,
		},
	}

	utils.SuccessResponse(ctx, "Article created", response)
}

// PUT /articles/:id
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid ID", err)
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Content     *string `json:"content"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	update := c.service.Client.Article.FindUnique(
		db.Article.ID.Equals(id),
	).With(db.Article.User.Fetch()).Update(
		db.Article.Title.SetIfPresent(input.Title),
		db.Article.Description.SetIfPresent(input.Description),
		db.Article.Content.SetIfPresent(input.Content),
	)

	article, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update article", err)
		return
	}

	var response models.ArticleResponse
	desc, ok := article.Description()
	if ok {
		desc = "(pas de description)"
	}
	response = models.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Description: desc,
		Content:     article.Content,
		CreatedAt:   article.CreateAt,
		UpdatedAt:   article.UpdateAt,
		User: &models.UserResponse{
			ID:        article.User().ID,
			Name:      article.User().Name,
			Email:     article.User().Email,
			Role:      string(article.User().Role),
			CreatedAt: article.User().CreateAt,
			UpdatedAt: article.User().UpdateAt,
		},
	}

	utils.SuccessResponse(ctx, "Article updated", response)
}

// DELETE /articles/:id
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid ID", err)
		return
	}

	_, err = c.service.Client.Article.FindUnique(
		db.Article.ID.Equals(id),
	).Delete().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to delete article", err)
		return
	}

	utils.SuccessResponse(ctx, "Article deleted", nil)
}
