package main

import (
	"fmt"
	"uupp-backend/database"
	"uupp-backend/handlers"
	"uupp-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {

	database.ConnectDB()
	app := fiber.New()

	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)
	app.Get("/protected", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("This is a protected route")
	})
	fmt.Println("Server is running on port 3006")
	app.Listen(":3006")
}
