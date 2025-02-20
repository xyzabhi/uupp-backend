package database

import (
	"log"
	"uupp-backend/config"
	"uupp-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.LoadConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db

	// Drop the users table if it exists
	if DB.Migrator().HasTable(&models.User{}) {
		if err := DB.Migrator().DropTable(&models.User{}); err != nil {
			log.Fatal("Failed to drop users table:", err)
		}
	}

	// Create the users table with the updated schema
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected successfully")
}
