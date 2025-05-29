package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/mihiryadav20/goauth/internal/auth"
	"github.com/mihiryadav20/goauth/internal/config"
	"github.com/mihiryadav20/goauth/internal/handlers"
	"github.com/mihiryadav20/goauth/internal/routes"
	"github.com/mihiryadav20/goauth/internal/trello"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Initialize configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error handling
			code := fiber.StatusInternalServerError

			// Check for specific error types
			e, ok := err.(*fiber.Error)
			if ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})


	// Use logger middleware
	app.Use(logger.New())


	// Initialize services
	trelloClient := trello.NewClient(cfg.TrelloAPIKey)
	authSvc := auth.NewService(cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, trelloClient, authSvc)

	// Setup routes
	routes.SetupRoutes(app, authHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
