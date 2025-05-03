package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/models"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.ServiceManager
}

func NewUserController(service *services.ServiceManager) *UserController {
	return &UserController{service: service}
}

// GET /users
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.service.Client.User.FindMany().Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to fetch users", err)
		return
	}
	var response []models.UserResponse

	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreateAt,
			UpdatedAt: user.UpdateAt,
		})
	}
	utils.SuccessResponse(ctx, "Users fetched successfully", response)
}

// GET /user/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid user ID", err)
		return
	}

	user, err := c.service.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(c.service.Ctx)
	if err != nil || user == nil {
		utils.ErrorResponse(ctx, 404, "User not found", err)
		return
	}

	var response = models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreateAt,
		UpdatedAt: user.UpdateAt,
	}

	utils.SuccessResponse(ctx, "User fetched successfully", response)
}

// GET /users/search?q=keyword
func (c *UserController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if strings.TrimSpace(query) == "" {
		utils.ErrorResponse(ctx, 400, "Search query is required", nil)
		return
	}

	users, err := c.service.Client.User.FindMany(
		db.User.Name.Contains(query),
	).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Search failed", err)
		return
	}

	var response []models.UserResponse

	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreateAt,
			UpdatedAt: user.UpdateAt,
		})
	}

	utils.SuccessResponse(ctx, "Search results", response)
}

// POST /user
func (c *UserController) CreateUser(ctx *gin.Context) {
	var input struct {
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required"`
		Password string  `json:"password" binding:"required"`
		Role     db.Role `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to hash password", err)
		return
	}

	user, err := c.service.Client.User.CreateOne(
		db.User.Name.Set(input.Name),
		db.User.Email.Set(input.Email),
		db.User.Password.Set(hashedPassword),
		db.User.Role.Set(input.Role),
	).Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to create user", err)
		return
	}

	var response = models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreateAt,
		UpdatedAt: user.UpdateAt,
	}
	utils.SuccessResponse(ctx, "User created", response)
}

// PUT /users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid user ID", err)
		return
	}

	var input struct {
		Name  *string  `json:"name" binding:"required"`
		Email *string  `json:"email" binding:"required"`
		Role  *db.Role `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid input", err)
		return
	}

	update := c.service.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.Name.SetIfPresent(input.Name),
		db.User.Email.SetIfPresent(input.Email),
		db.User.Role.SetIfPresent(input.Role),
	)

	user, err := update.Exec(c.service.Ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to update user", err)
		return
	}

	var response = models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreateAt,
		UpdatedAt: user.UpdateAt,
	}
	utils.SuccessResponse(ctx, "User updated successfully", response)
}

// DELETE /users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid user ID", err)
		return
	}

	_, err = c.service.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(c.service.Ctx)

	if err != nil {
		utils.ErrorResponse(ctx, 500, "Failed to delete user", err)
		return
	}

	utils.SuccessResponse(ctx, "User deleted successfully", nil)
}
