package main

import (
	"bytecrate/internal"
	"bytecrate/internal/database"
	"log"
	"os"
)

func main() {
	database.Connect()
	router := internal.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Backend running on :8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
