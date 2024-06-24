package main

import (
	"log"

	"github.com/ajay/exptest/config"
	"github.com/ajay/exptest/handlers"
	"github.com/ajay/exptest/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001,http://localhost:3002,http://localhost:3003,http://localhost:3000,http://localhost:8080",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Serve static files
	app.Static("/", "./public")

	// Initialize MongoDB client
	config.ConnectDB()

	// Public routes
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.Register)
	app.Post("/logout", handlers.Logout)

	// Protected routes
	app.Use(middleware.AuthMiddleware)

	// Transaction routes
	app.Get("/transactions", handlers.GetTransactions)
	app.Post("/transactions", handlers.CreateTransaction)
	app.Put("/transactions/:id", handlers.UpdateTransaction)
	app.Delete("/transactions/:id", handlers.DeleteTransaction)

	log.Fatal(app.Listen(":5000"))
}
