package main

import (
	"bytecrate/internal"
	"bytecrate/internal/database"
	"log"
	"os"
)

func main() {
	db := database.Connect()
	router := internal.NewRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Backend running on :8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
