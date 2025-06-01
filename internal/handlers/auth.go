package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mihiryadav20/goauth/internal/auth"
	"github.com/mihiryadav20/goauth/internal/config"
	"github.com/mihiryadav20/goauth/internal/models"
	"github.com/mihiryadav20/goauth/internal/trello"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	config     *config.Config
	trello     *trello.Client
	AuthSvc    *auth.Service
	stateStore map[string]time.Time
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(cfg *config.Config, trelloClient *trello.Client, authSvc *auth.Service) *AuthHandler {
	return &AuthHandler{
		config:     cfg,
		trello:     trelloClient,
		AuthSvc:    authSvc,
		stateStore: make(map[string]time.Time),
	}
}

// InitiateAuth handles the OAuth initiation
func (h *AuthHandler) InitiateAuth(c *fiber.Ctx) error {
	state, err := h.AuthSvc.GenerateState()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate the state",
		})
	}

	// Store state with expiration
	h.stateStore[state] = time.Now().Add(time.Minute * 10)

	// Build Trello authorization URL
	authURL := fmt.Sprintf(
		"https://trello.com/1/authorize?expiration=30days&name=GoAuth&scope=read&response_type=token&key=%s&return_url=%s&callback_method=fragment&state=%s",
		h.config.TrelloAPIKey,
		url.QueryEscape(h.config.FrontendURL+"/callback"),
		state,
	)

	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"state":    state,
	})
}

// HandleCallback handles the OAuth callback
func (h *AuthHandler) HandleCallback(c *fiber.Ctx) error {
	// Get query parameters
	code := c.Query("code")
	state := c.Query("state")

	// Validate state
	expiry, exists := h.stateStore[state]
	if !exists || time.Now().After(expiry) {
		delete(h.stateStore, state)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or expired state",
		})
	}

	// Clean up used state
	delete(h.stateStore, state)

	// If code is missing, return error
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing authorization code",
		})
	}

	// Exchange code for token
	token, err := h.trello.ExchangeCodeForToken(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange code for token",
		})
	}

	// Get user info from Trello
	userInfo, err := h.trello.GetUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info from Trello",
		})
	}

	// Generate JWT
	jwtToken, err := h.AuthSvc.GenerateJWT(userInfo.ID, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT",
		})
	}

	return c.JSON(models.AuthResponse{
		Token:    jwtToken,
		UserInfo: userInfo,
	})
}

// GetProfile returns the authenticated user's profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	trelloToken := c.Locals("trello_token").(string)

	// Get user info from Trello
	userInfo, err := h.trello.GetUserInfo(trelloToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info from Trello",
		})
	}

	return c.JSON(fiber.Map{
		"user": userInfo,
	})
}

// GetBoards returns the authenticated user's Trello boards
func (h *AuthHandler) GetBoards(c *fiber.Ctx) error {
	trelloToken := c.Locals("trello_token").(string)

	// Get boards from Trello
	boards, err := h.trello.GetUserBoards(trelloToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get boards from Trello",
		})
	}

	return c.JSON(fiber.Map{
		"boards": boards,
	})
}

// RegisterTokenRequest represents the request body for token registration
type RegisterTokenRequest struct {
	TrelloToken string `json:"trello_token"`
}

// RegisterToken handles registering a Trello token and generating a JWT
func (h *AuthHandler) RegisterToken(c *fiber.Ctx) error {
	// Parse request body
	req := new(RegisterTokenRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate token
	if req.TrelloToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing Trello token",
		})
	}

	// Get user info from Trello to validate the token
	userInfo, err := h.trello.GetUserInfo(req.TrelloToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Trello token",
		})
	}

	// Generate JWT
	jwtToken, err := h.AuthSvc.GenerateJWT(userInfo.ID, req.TrelloToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT",
		})
	}

	return c.JSON(models.AuthResponse{
		Token:    jwtToken,
		UserInfo: userInfo,
	})
}
