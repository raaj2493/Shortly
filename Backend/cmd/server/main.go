package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/raaj2493/Shortly/Backend/internals/config"
	"github.com/raaj2493/Shortly/Backend/internals/database"
	"github.com/raaj2493/Shortly/Backend/internals/handlers"
	"github.com/raaj2493/Shortly/Backend/internals/middleware"
	"github.com/raaj2493/Shortly/Backend/internals/repository"
	"github.com/raaj2493/Shortly/Backend/internals/services"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Connect to PostgreSQL
	repository.InitDB(cfg.DatabaseDSN)
	db := repository.DB

	// 3. Connect to Redis
	repository.InitRedis(cfg.RedisAddr, cfg.RedisPass)
	redisClient := repository.RedisClient

	// 4. Run database migrations
	database.RunMigration(cfg.DatabaseDSN, "internals/migration/migrations")

	// 5. Create repository instances
	userRepo := repository.NewUserRepository(db)

	// 6. Create service instances
	authService := services.NewAuthService(cfg, redisClient, userRepo)

	// 7. Create handler instances
	authHandler := handlers.NewAuthHandler(authService)
	urlHandler := middleware.NewURLHandler()

	// 8. Set up HTTP router using standard mux
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Protected routes – we'll wrap them with the auth middleware
	sessionChecker := func(userID int) (bool, error) {
		exists, err := redisClient.Exists(context.Background(), "session:"+strconv.Itoa(userID)).Result()
		return exists == 1, err
	}
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret, sessionChecker)

	mux.HandleFunc("/api/urls", func(w http.ResponseWriter, r *http.Request) {
		protectedCreateURL := authMiddleware(http.HandlerFunc(urlHandler.CreateShortURL))
		protectedCreateURL.ServeHTTP(w, r)
	})

	// Public redirect route (handled by URL handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only handle GET requests; ignore other routes like /health etc.
		if r.Method == http.MethodGet && r.URL.Path != "/health" {
			urlHandler.RedirectShortURL(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// 9. Start HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Run server in a goroutine so we can listen for shutdown signals
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// 10. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// healthHandler returns a simple health check.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}