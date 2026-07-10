package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raaj2493/Shortly/Backend/config"
	"github.com/raaj2493/Shortly/Backend/database"
	"github.com/raaj2493/Shortly/Backend/handlers"
	"github.com/raaj2493/Shortly/Backend/logger"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize structured logger
	log := logger.New(cfg.LogLevel, cfg.LogFormat)
	log.Info("starting url-shortener",
		"server_port", cfg.ServerPort,
		"db_host", cfg.DBHost,
	)

	// 3. Connect to PostgreSQL
	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close() // close the pool on exit
	log.Info("connected to PostgreSQL")

	// 4. Run database migrations
	if err := database.RunMigrations(cfg); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("database migrations applied")

	// 5. Set up health check HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// 6. Start HTTP server in a goroutine
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Channel to listen for OS signals (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in background
	go func() {
		log.Info("starting health server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("health server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Wait for shutdown signal
	sig := <-quit
	log.Info("received signal, shutting down", "signal", sig.String())

	// Give outstanding requests a deadline to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server exited gracefully")
}