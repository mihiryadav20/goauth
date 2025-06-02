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
	CallbackURL       string
	FrontendURL       string
	JWTSecret         string
	SessionExpiration time.Duration
	// Notion OAuth credentials
	NotionClientID     string
	NotionClientSecret string
	NotionRedirectURI  string
	NotionAuthURL      string
	NotionTokenURL     string
	NotionAPIBaseURL   string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		TrelloAPIKey:      getEnv("TRELLO_API_KEY", ""),
		TrelloAPISecret:   getEnv("TRELLO_API_SECRET", ""),
		CallbackURL:       getEnv("CALLBACK_URL", "http://localhost:3000/auth/trello/callback"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3001"),
		JWTSecret:         getEnv("JWT_SECRET", "your-jwt-secret-key"),
		SessionExpiration: time.Hour * 24, // 24 hours
		// Notion OAuth credentials
		NotionClientID:     getEnv("NOTION_CLIENT_ID", ""),
		NotionClientSecret: getEnv("NOTION_CLIENT_SECRET", ""),
		NotionRedirectURI:  getEnv("NOTION_REDIRECT_URI", "http://localhost:3000/auth/notion/callback"),
		NotionAuthURL:      getEnv("NOTION_AUTH_URL", "https://api.notion.com/v1/oauth/authorize"),
		NotionTokenURL:     getEnv("NOTION_TOKEN_URL", "https://api.notion.com/v1/oauth/token"),
		NotionAPIBaseURL:   getEnv("NOTION_API_BASE_URL", "https://api.notion.com/v1"),
	}
}

// Validate checks if required configuration is set
func (c *Config) Validate() error {
	// We don't make Notion credentials required yet since we're just adding them
	if c.TrelloAPIKey == "" || c.TrelloAPISecret == "" {
		return fmt.Errorf("TRELLO_API_KEY and TRELLO_API_SECRET must be set")
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
