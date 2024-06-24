package handlers

import (
	"github.com/ajay/exptest/config"
	"github.com/ajay/exptest/models"
	"github.com/ajay/exptest/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var request RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate if email already exists
	existingEmail := &models.User{}
	err := config.Db.Collection("users").FindOne(c.Context(), bson.M{"email": request.Email}).Decode(existingEmail)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Validate if username already exists
	existingUsername := &models.User{}
	err = config.Db.Collection("users").FindOne(c.Context(), bson.M{"username": request.Username}).Decode(existingUsername)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create new user object
	newUser := models.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Email:    request.Email,
	}

	// Insert user into database
	_, err = config.Db.Collection("users").InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Return success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// Login handles user login
func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user by username
	var user models.User
	err := config.Db.Collection("users").FindOne(c.Context(), bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return token
	return c.JSON(fiber.Map{
		"token": token,
	})
}

// Logout handles user logout (if needed)
func Logout(c *fiber.Ctx) error {
	// Typically, for JWT, no server-side action is required to log out.
	// Here we can simply return a message indicating logout.

	// Optionally, you could add token blacklisting or other logic here.
	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
