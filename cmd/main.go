package main

import (
	"fmt"
	"uupp-backend/database"
	"uupp-backend/handlers"
	"uupp-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()
	app := fiber.New()

	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)
	app.Get("/protected", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("This is a protected route")
	})
	fmt.Println("Server is running on port 3005")
	app.Listen(":3005")
}
