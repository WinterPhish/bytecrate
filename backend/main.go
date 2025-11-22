package main

import (
	"bytecrate/internal"
	"log"
)

func main() {
	router := internal.NewRouter()

	log.Println("Backend running on :8080")
	router.Run(":8080")
}
