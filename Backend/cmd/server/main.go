package main

import(
	"fmt"

	"github.com/raaj2493/Shortly/Backend/internal/config"
	
)

func main(){
	// 1. Load the configuration
	cfg := config.LoadConfig()

	// 2. Start the server
	fmt.Printf("🚀 Starting URL Shortener API in %s mode on port %s...\n", cfg.Env, cfg.Port)
}

