package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	service *services.ServiceManager
}

func NewArticleController(service *services.ServiceManager) *ArticleController {
	return &ArticleController{service: service}
}

// GET /articles
func (c *ArticleController) GetArticles(ctx *gin.Context) {
	articles, err := c.service.Client.Article.FindMany().Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to fetch articles", err)
		return
	}
	utils.SuccessResponse(ctx, "Articles fetched successfully", articles)
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
	).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "Article not found", err)
		return
	}

	utils.SuccessResponse(ctx, "Article fetched successfully", article)
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
	).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Search failed", err)
		return
	}

	utils.SuccessResponse(ctx, "Search results", articles)
}

// POST /article
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Content     string `json:"content" binding:"required"`
		UserID      int    `json:"userId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	article, err := c.service.Client.Article.CreateOne(
		db.Article.Title.Set(input.Title),
		db.Article.Content.Set(input.Content),
		db.Article.User.Link(
			db.User.ID.Equals(input.UserID),
		),
		db.Article.Description.Set(input.Description),
	).Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to create article", err)
		return
	}

	utils.SuccessResponse(ctx, "Article created", article)
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
	).Update(
		db.Article.Title.SetIfPresent(input.Title),
		db.Article.Description.SetIfPresent(input.Description),
		db.Article.Content.SetIfPresent(input.Content),
	)

	article, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update article", err)
		return
	}

	utils.SuccessResponse(ctx, "Article updated", article)
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
