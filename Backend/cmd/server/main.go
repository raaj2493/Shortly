package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raaj2493/Shortly/Backend/internals/config"
	"github.com/raaj2493/Shortly/Backend/internals/database"
	"github.com/raaj2493/Shortly/Backend/internals/handler"
	"github.com/raaj2493/Shortly/Backend/internals/repository"
	"github.com/raaj2493/Shortly/Backend/internals/services"
)

func main() {

	//1. Load the Config
	cfg := config.Load()

	//2. Initialize the DB
	db := repository.NewDB(cfg.DatabaseURL)
	defer db.Close()

	//3. Initialize the Redis Client
	redisClient := repository.NewRedisClient(cfg.RedisURL)
	defer redisClient.Close()

	//4. Running DB migration
	database.RunMigrations(cfg.DatabaseURL)


	//Initialise Layers
	urlRepository := repository.NewURLRepository(db)
	urlService := services.NewURLService(urlRepository, redisClient)
	urlHandler := handler.NewURLHandler(urlService)


	//5. Setting UP the server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})


	//6. Initialize the Services
      addr := ":" + cfg.ServerPort
	  log.Printf("Starting server on %s", addr)

	  err := http.ListenAndServe(addr, mux)
	  if err != nil {
		log.Fatal(err)
	  }


// URL shortener endpoints
	mux.HandleFunc("POST /api/urls", urlHandler.CreateShortURL)
	mux.HandleFunc("GET /{shortCode}", urlHandler.RedirectToOriginal)

	// 6. HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	// Run server in a goroutine so we can listen for shutdown signals
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 7. Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 8. Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
