// Package handlers contains HTTP request handlers
package handlers

// Import necessary packages for authentication, database operations, and web functionality
import (
	"time"
	"uupp-backend/database"
	"uupp-backend/models"

	"github.com/gofiber/fiber/v2"  // Web framework
	"github.com/golang-jwt/jwt/v5" // JWT token handling
	"golang.org/x/crypto/bcrypt"   // Password hashing
)

// JWT secret key used for signing and verifying tokens
var jwtSecret = []byte("your-256-bit-secret")

// RegisterUser handles new user registration
func RegisterUser(c *fiber.Ctx) error {
	// Create a new User struct to store the registration data
	var user models.User

	// Parse the request body into the user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	// Store the hashed password in the user struct
	user.Password = string(hashedPassword)

	// Attempt to create the user in the database
	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Return success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// LoginUser handles user authentication and token generation
func LoginUser(c *fiber.Ctx) error {
	// Create a struct to store the login credentials
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Look up the user in the database by username
	var user models.User
	if result := database.DB.Where("username = ?", input.Username).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	// Create a new JWT token with user claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return the JWT token to the client
	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
