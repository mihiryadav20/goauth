package models

import "time"

// User represents a user in our system
type User struct {
	ID          string    `json:"id"`
	TrelloToken string    `json:"trello_token"`
	CreatedAt   time.Time `json:"created_at"`
}

// TrelloUserInfo represents basic user info from Trello
type TrelloUserInfo struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}
