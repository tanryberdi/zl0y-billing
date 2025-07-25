package main

import (
	"log"

	"zl0y-billing/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// print the configuration for debugging purposes
	log.Printf("Configuration: %+v\n", cfg)

}
