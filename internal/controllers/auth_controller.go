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
	var response models.LoginResponse
	response.Token = token
	response.User = models.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Email:     userRecord.Email,
		Role:      string(userRecord.Role),
		CreatedAt: userRecord.CreateAt,
		UpdatedAt: userRecord.UpdateAt,
	}
	utils.SuccessResponse(ctx, "Login successful", response)

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

	// CHeck if exist users if exist, then set role to user, else set to admin
	var role db.Role
	existingUsers, err := c.service.Client.User.FindMany().Exec(c.service.Ctx)
	if err != nil {
		log.Println("Error checking existing users:", err)
		utils.InternalServerErrorResponse(ctx, "Failed to check existing users")
		return
	}
	if len(existingUsers) > 0 {
		role = db.RoleUser // Set role to user if users already exist
	} else {
		role = db.RoleAdmin // Set role to admin if no users exist
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
		db.User.Role.Set(role),
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

	var response models.LoginResponse
	response.Token = token
	response.User = models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      string(newUser.Role),
		CreatedAt: newUser.CreateAt,
		UpdatedAt: newUser.UpdateAt,
	}
	utils.CreatedResponse(ctx, "User registered successfully", response)
	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
}
func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	utils.SuccessResponse(ctx, "Logged out successfully", nil)
}
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	// Implement reset password logic here
}
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var input struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
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

	// Rechercher l'utilisateur
	userRecord, err := c.service.Client.User.FindFirst(db.User.Email.Equals(userClaims.Email)).Exec(c.service.Ctx)
	if err != nil || userRecord == nil {
		utils.BadRequestResponse(ctx, "User not found")
		return
	}

	// Vérifier l'ancien mot de passe
	if !utils.CheckPasswordHash(input.OldPassword, userRecord.Password) {
		utils.BadRequestResponse(ctx, "Old password is incorrect")
		return
	}

	// Hasher le nouveau mot de passe
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to hash new password")
		return
	}

	// Mettre à jour le mot de passe
	_, err = c.service.Client.User.FindUnique(db.User.ID.Equals(userRecord.ID)).Update(
		db.User.Password.Set(hashedPassword),
	).Exec(c.service.Ctx)

	if err != nil {
		utils.InternalServerErrorResponse(ctx, "Failed to update password")
		return
	}

	utils.SuccessResponse(ctx, "Password changed successfully", nil)
}
