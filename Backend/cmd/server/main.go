package main

import (
	"fmt"
	"log"

	"github.com/raaj2493/Shortly/Backend/internal/config"
	"github.com/raaj2493/Shortly/Backend/internal/db"
	"github.com/raaj2493/Shortly/Backend/internal/repository"
)

func main() {
	// 1. Run the checklist
	cfg := config.LoadConfig()
	fmt.Printf("🚀 Starting URL Shortener API in %s mode on port %s...\n", cfg.Env, cfg.Port)

	// 2. Connect to the pipeline bucket
	databasePool := db.InitDB(cfg)
	defer databasePool.Close()

	// 3. Initialize our Repository Robots
	userRepo := repository.NewUserRepository(databasePool)
	_ = userRepo // Suppresses the "unused variable" error until we build handlers next!

	log.Println("🤖 Repository layer initialized and waiting for instructions...")
}