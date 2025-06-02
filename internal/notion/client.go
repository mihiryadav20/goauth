package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mihiryadav20/goauth/internal/config"
	"github.com/mihiryadav20/goauth/internal/models"
)

// Client represents a Notion API client
type Client struct {
	ClientID     string
	ClientSecret string
	Config       *config.Config
}

// NewClient creates a new Notion client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		ClientID:     cfg.NotionClientID,
		ClientSecret: cfg.NotionClientSecret,
		Config:       cfg,
	}
}

// GetAuthorizationURL returns the URL to redirect users to for authorization
func (c *Client) GetAuthorizationURL(state string) string {
	redirectURI := c.Config.FrontendURL + "/auth/notion/callback"
	return fmt.Sprintf(
		"https://api.notion.com/v1/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		url.QueryEscape(c.ClientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
}

// TokenResponse represents the response from Notion's OAuth token endpoint
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	BotID       string `json:"bot_id"`
	Workspace   struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"workspace"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceIcon string `json:"workspace_icon"`
	Owner         struct {
		Type string `json:"type"`
		User struct {
			Object    string `json:"object"`
			ID        string `json:"id"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
	} `json:"owner"`
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (c *Client) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	redirectURI := c.Config.FrontendURL + "/auth/notion/callback"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://api.notion.com/v1/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code for token: %s - %s", resp.Status, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo gets user information from Notion
func (c *Client) GetUserInfo(accessToken string) (*models.NotionUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.notion.com/v1/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s - %s", resp.Status, string(body))
	}

	var userInfo models.NotionUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
