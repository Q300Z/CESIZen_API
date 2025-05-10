package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse envoie une réponse 200 avec un message et des données
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse envoie une réponse 201
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse envoie une réponse d’erreur personnalisée
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	c.JSON(statusCode, ApiResponse{
		Success: false,
		Message: message,
		Error:   errMsg,
	})
}

// NotFoundResponse envoie une réponse 404
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ApiResponse{
		Success: false,
		Message: message,
	})
}

// BadRequestResponse envoie une réponse 400
func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ApiResponse{
		Success: false,
		Message: message,
	})
}

// UnauthorizedResponse envoie une réponse 401
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ApiResponse{
		Success: false,
		Message: message,
	})
}

// ForbiddenResponse envoie une réponse 403
func ForbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ApiResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerErrorResponse envoie une réponse 500
func InternalServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ApiResponse{
		Success: false,
		Message: message,
	})
}
