package main

import (
	"bytecrate/internal/http"
	"log"
)

func main() {
	router := http.NewRouter()

	log.Println("Backend running on :8080")
	router.Run(":8080")
}
