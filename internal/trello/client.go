package trello

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mihiryadav20/goauth/internal/config"
	"github.com/mihiryadav20/goauth/internal/models"
)

// Client represents a Trello API client
type Client struct {
	APIKey    string
	APISecret string
	Config    *config.Config
}

// NewClient creates a new Trello client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:    cfg.TrelloAPIKey,
		APISecret: cfg.TrelloAPISecret,
		Config:    cfg,
	}
}

// GetUserInfo gets user information from Trello
// GetUserInfo gets user information from Trello
func (c *Client) GetUserInfo(token string) (*models.TrelloUserInfo, error) {
	// Build request URL
	requestURL := fmt.Sprintf(
		"https://api.trello.com/1/members/me?key=%s&token=%s",
		c.APIKey,
		token,
	)

	// Make request to Trello API
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trello API returned status %d", resp.StatusCode)
	}

	// Decode response
	var userInfo models.TrelloUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// ExchangeCodeForToken exchanges an authorization code for a Trello token
func (c *Client) ExchangeCodeForToken(code string) (string, error) {
	// Trello doesn't have a standard OAuth token endpoint
	// Instead, we need to use the code directly as a token
	// This is a simplification of the process - in a real app, you might want to validate the code
	// by making a test request to the Trello API
	
	// Build request URL for validation
	requestURL := fmt.Sprintf(
		"https://api.trello.com/1/members/me?key=%s&token=%s",
		c.APIKey,
		code,
	)

	// Make request to Trello API to validate the token
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to validate Trello token, API returned status %d", resp.StatusCode)
	}

	// If we got here, the code is valid and can be used as a token
	return code, nil
}

// GetUserBoards gets all boards for the authenticated user
func (c *Client) GetUserBoards(token string) ([]map[string]interface{}, error) {
	// Build request URL
	requestURL := fmt.Sprintf(
		"https://api.trello.com/1/members/me/boards?key=%s&token=%s",
		c.APIKey,
		token,
	)

	// Make request to Trello API
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trello API returned status %d", resp.StatusCode)
	}

	// Decode response
	var boards []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&boards); err != nil {
		return nil, err
	}

	return boards, nil
}
