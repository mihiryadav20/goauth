package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	TrelloAPIKey      string
	TrelloAPISecret   string
	NotionClientID    string
	NotionClientSecret string
	CallbackURL       string
	FrontendURL       string
	JWTSecret         string
	SessionExpiration time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		TrelloAPIKey:      getEnv("TRELLO_API_KEY", ""),
		TrelloAPISecret:   getEnv("TRELLO_API_SECRET", ""),
		NotionClientID:    getEnv("NOTION_CLIENT_ID", ""),
		NotionClientSecret: getEnv("NOTION_CLIENT_SECRET", ""),
		CallbackURL:       getEnv("CALLBACK_URL", "http://localhost:3000/auth/trello/callback"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3001"),
		JWTSecret:         getEnv("JWT_SECRET", "your-jwt-secret-key"),
		SessionExpiration: time.Hour * 24, // 24 hours
	}
}

// Validate checks if required configuration is set
func (c *Config) Validate() error {
	if c.TrelloAPIKey == "" || c.TrelloAPISecret == "" {
		return fmt.Errorf("TRELLO_API_KEY and TRELLO_API_SECRET must be set")
	}
	if c.NotionClientID == "" || c.NotionClientSecret == "" {
		return fmt.Errorf("NOTION_CLIENT_ID and NOTION_CLIENT_SECRET must be set")
	}
	return nil
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
