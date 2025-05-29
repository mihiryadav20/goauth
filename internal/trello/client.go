package trello

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mihiryadav20/goauth/internal/models"
)

// Client represents a Trello API client
type Client struct {
	APIKey string
}

// NewClient creates a new Trello client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
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
