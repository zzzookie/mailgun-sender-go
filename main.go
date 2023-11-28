package main

import (
	"os"

	"mailgun-sender-go/models"
	"mailgun-sender-go/storage"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		panic("Could not load the database")
	}

	err = models.Migrate(db)
	if err != nil {
		panic("Could not migrate DB")
	}

	// r := Repository{
	// 	DB: db,
	// }
}
