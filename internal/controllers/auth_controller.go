package controllers

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/models"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.ServiceManager
}

func NewAuthController(service *services.ServiceManager) *AuthController {
	return &AuthController{service: service}
}

// Login
func (c *AuthController) Login(ctx *gin.Context) {
	var user struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	// Validate username and password
	if user.Email == "" || user.Password == "" {
		utils.BadRequestResponse(ctx, "Email and password are required")
		return
	}
	// Check if user exists and password is correct
	userRecord, err := c.service.Client.User.FindFirst(db.User.Email.Equals(user.Email)).Exec(c.service.Ctx)
	if err != nil || userRecord == nil {
		utils.BadRequestResponse(ctx, "Invalid email or password")
		return
	}

	// Check if password matches
	if !utils.CheckPasswordHash(user.Password, userRecord.Password) {
		utils.BadRequestResponse(ctx, "Invalid email or password")
		return
	}
	// Generate JWT token
	token, err := utils.GenerateJWT(userRecord)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to generate token")
		return
	}
	// Return token to client
	utils.SuccessResponse(ctx, "Login successful", gin.H{"token": token})

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
}
func (c *AuthController) Register(ctx *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	// Validate username and password
	if user.Email == "" || user.Password == "" || user.Username == "" {
		utils.BadRequestResponse(ctx, "Email and password are required")
		return
	}
	// Check if user already exists
	existingUser, err := c.service.Client.User.FindFirst(db.User.Email.Equals(user.Email)).Exec(c.service.Ctx)
	if err != nil {
		if err.Error() != "ErrNotFound" {
			log.Println("Error checking for existing user:", err)
			utils.InternalServerErrorResponse(ctx, "Failed to check for existing user")
			return
		}
	} else if existingUser != nil {
		utils.BadRequestResponse(ctx, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to hash password")
		return
	}

	// Create new user
	newUser, err := c.service.Client.User.CreateOne(
		db.User.Name.Set(user.Username),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(hashedPassword),
	).Exec(c.service.Ctx)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(newUser)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to generate token")
		return
	}

	var response struct {
		Token string              `json:"token"`
		User  models.UserResponse `json:"user"`
	}
	response.Token = token
	// Exclude password from the response
	newUser.Password = ""
	response.User = models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      string(newUser.Role),
		CreatedAt: newUser.CreateAt,
		UpdatedAt: newUser.UpdateAt,
	}
	utils.SuccessResponse(ctx, "User registered successfully", response)
	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
}
func (c *AuthController) Logout(ctx *gin.Context) {
	// Implement logout logic here
}
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	// Implement reset password logic here
}
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// Implement change password logic here
}
