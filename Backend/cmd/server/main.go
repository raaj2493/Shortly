package main

import (
	"log"
	"net/http"

	"github.com/raaj2493/Shortly/Backend/internals/config"
	"github.com/raaj2493/Shortly/Backend/internals/handler"
	"github.com/raaj2493/Shortly/Backend/internals/repository"
	"github.com/raaj2493/Shortly/Backend/internals/service"
	"github.com/raaj2493/Shortly/Backend/internals/database"
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

}
