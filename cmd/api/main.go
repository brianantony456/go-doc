package main

import (
	"log"

	"github.com/brianantony456/go-doc/internal/router"
)

func main() {
	r, _, err := router.SetupRouter()
	if err != nil {
		log.Fatal("Failed to set up router and connect to database:", err)
	}

	// Defer closing the database connection if necessary (depending on the database driver)
	// defer db.Close()

	if err := r.Run("localhost:8000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
