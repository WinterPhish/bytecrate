package database

import (
	"bytecrate/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=postgres user=postgres password=postgres dbname=bytecrate port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.User{})

	return db
}