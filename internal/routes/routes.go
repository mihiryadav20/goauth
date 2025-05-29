package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mihiryadav20/goauth/internal/auth"
	"github.com/mihiryadav20/goauth/internal/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	authSvc := authHandler.AuthSvc
	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Trello OAuth API",
		})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Get("/trello", authHandler.InitiateAuth)
	auth.Get("/trello/callback", authHandler.HandleCallback)

	// Protected routes
	api := app.Group("/api", authMiddleware(authSvc))
	api.Get("/profile", authHandler.GetProfile)
	api.Get("/boards", authHandler.GetBoards)
}

// authMiddleware verifies the JWT token and sets user info in the context
func authMiddleware(authSvc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Parse token (remove 'Bearer ' prefix)
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		claims, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user info in context
		if userID, ok := claims["user_id"].(string); ok {
			c.Locals("user_id", userID)
		}

		if trelloToken, ok := claims["trello_token"].(string); ok {
			c.Locals("trello_token", trelloToken)
		}

		return c.Next()
	}
}
