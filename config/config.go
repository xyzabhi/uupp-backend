package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() string {
	if err := godotenv.Load(); err != nil {
		// log.Fatal("Error loading .env file")
		return "host=localhost user=postgres password=Green@8083 dbname=uupp_db port=5432 sslmode=disable"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}
