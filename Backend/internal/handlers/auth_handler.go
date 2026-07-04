package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/Shortly/Backend/internal/models"
	"github.com/raaj2493/Shortly/Backend/internal/repository"
)

type AuthHandler struct {
	repo *repository.UserRepository
}
func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// RegisterInput defines the exact JSON structure we expect from the client
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterUser(c *gin.Context){
	var input RegisterInput

	// 1. Automatically bind incoming JSON and validate required fields!
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required and email must be valid"})
		return
	}

	// 2. Map inputs to our database user model
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: input.Password, 
	}

	// 3. Send data to our repository robot using the request context
	err := h.UserRepository.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}

	// 4. Return a success status along with the clean user model data
	c.JSON(http.StatusCreated, user)
}