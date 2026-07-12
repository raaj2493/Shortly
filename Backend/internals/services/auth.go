package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/raaj2493/Shortly/Backend/internals/config"
	"github.com/raaj2493/Shortly/Backend/internals/models"
	"github.com/raaj2493/Shortly/Backend/internals/repository"
)

type AuthService struct {
	config *config.Config
	redis  *redis.Client
	userRepo   *repository.UserRepository
}

func NewAuthService (config *config.Config , redis *redis.Client , userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		config: config,
		redis:  redis,
		userRepo: userRepo,
	}
}

func ( s *AuthService) Register (ctx context.Context , email , password string) (*models.User , error) {

	// Hash the password
	hashed , err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		Email: email,
		Password: string(hashed),
	}

	err = s.userRepo.CreateUser(ctx , user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// Remove password before returning (just in case)
	user.Password = ""
	return user, nil
}

// login service

func (s *AuthService) Login(ctx context.Context, email , password string) (string , error ){
	// find user by email 
	user , err := s.userRepo.GetUserByEmail(ctx , email)
	if err != nil {
		return "", fmt.Errorf("finding user: %w", err)
	}

	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// compare password with hash 
	err = bcrypt.CompareHashAndPassword([]byte(user.Password) , []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , jwt.MapClaims{
		"sub": strconv.Itoa(user.ID),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	// Store token in Redis with key "session:<userID>" and 24h TTL

	key := fmt.Sprintf("session:%d" , user.ID)
	err = s.redis.Set(ctx, key, tokenString, 24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("storing session: %w", err)
	}

	return tokenString, nil
}


