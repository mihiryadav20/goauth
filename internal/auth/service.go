package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mihiryadav20/goauth/internal/config"
)

// Service handles authentication related operations
type Service struct {
	jwtSecret         string
	sessionExpiration time.Duration
}

// NewService creates a new authentication service
func NewService(cfg *config.Config) *Service {
	return &Service{
		jwtSecret:         cfg.JWTSecret,
		sessionExpiration: cfg.SessionExpiration,
	}
}

// GenerateState generates a random state string for OAuth
func (s *Service) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateJWT generates a JWT token for the user
func (s *Service) GenerateJWT(userID, trelloToken string) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      userID,
		"trello_token": trelloToken,
		"exp":          time.Now().Add(s.sessionExpiration).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKeyType
}
