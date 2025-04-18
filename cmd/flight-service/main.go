package main

import (
	"log"

	"github.com/fehepe/flight-price-service/internal/config"
	"github.com/fehepe/flight-price-service/internal/server"
)

func main() {
	config.LoadEnv()

	port := config.Get("PORT", "3000")

	// Run the server
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
