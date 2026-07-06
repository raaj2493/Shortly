package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/raaj2493/Shortly/Backend/internal/config"
	"github.com/raaj2493/Shortly/Backend/internal/database"
	"github.com/raaj2493/Shortly/Backend/internal/handlers"
	"github.com/raaj2493/Shortly/Backend/internal/repository"
	"github.com/raaj2493/Shortly/Backend/internal/routes"
)

func main() {
	// 1. Run the config safety checklist (.env parsing)
	cfg := config.LoadConfig()

	// 2. Set Gin mode based on our configuration environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Connect to our PostgreSQL connection pool
	databasePool := database.InitDB(cfg)
	
	// Ensure the database connection pool closes cleanly when the server shuts down
	defer databasePool.Close()

	// 4. Initialize our Repository Robots (Database interaction layer)
	userRepo := repository.NewUserRepository(databasePool)

	// 5. Initialize our Web Handler Layer (Controllers)
	authHandler := handlers.NewAuthHandler(userRepo)

	// 6. Initialize the Gin Engine with default logger and crash recovery middleware
	r := gin.Default()

	// 7. Delegate all API network routes to our dedicated routes package
	routes.SetupRouters(r, authHandler)

	// 8. Fire up the Gin network server engine!
	serverAddress := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("🚀 Starting URL Shortener API in %s mode on port %s...\n", cfg.Env, cfg.Port)
	
	err := r.Run(serverAddress)
	if err != nil {
		log.Fatalf("Network server crashed unexpectedly: %v", err)
	}
}