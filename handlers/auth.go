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
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Validate all required fields
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check for unique username
	if user.IsUsernameExists(database.DB) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Check for unique email
	if user.IsEmailExists(database.DB) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user.Password = string(hashedPassword)

	// Create the user
	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

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

	// Create a user response without the password
	userResponse := fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"dob":       user.DOB,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	// Return both the JWT token and user data
	return c.JSON(fiber.Map{
		"token": tokenString,
		"user":  userResponse,
	})
}
