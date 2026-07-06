package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/Shortly/Backend/internal/models"
	"github.com/raaj2493/Shortly/Backend/internal/repository"
	"github.com/raaj2493/Shortly/Backend/internal/utils"
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

	// 2. Hash the raw password before it ever touches the database layer
	hashedPassword, err := utils.HashPass(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to securely process password"})
		return
	}

	// 2. Map inputs to our database user model
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword, 
	}

	// 3. Send data to our repository robot using the request context
	err = h.repo.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}

	// 4. Return a success status along with the clean user model data
	c.JSON(http.StatusCreated, user)
}


// LoginInput defines the structure for incoming login credentials
// LoginInput defines the structure for incoming login credentials
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginUser verifies email/password and returns a fresh JWT token
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var input LoginInput

	// 1. Bind and validate incoming JSON fields
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email and password are required"})
		return
	}

	// 2. Fetch the user from the database by email
	user, err := h.repo.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		// Security tip: Don't tell hackers if the email was wrong or password was wrong. Keep it vague!
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. Compare the raw password with our secure database blender hash
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 4. Retrieve our JWT secret string from application context config
	// For now, we will grab it from environment variables using a quick helper, or pass it into our handler initialization later.
	// Let's generate the token!
	// (Note: To keep this clean, we pull the config out)
	// For simplicity in this step, we can grab it directly via our utility or inject it.
	secret := "super_secret_shortly_key_2026_dont_leak_this" // We will clean this injection up on refactor!

	token, err := utils.GenerateToken(user.ID, secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate login session"})
		return
	}

	// 5. Send back the token!
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful!",
	})
}

