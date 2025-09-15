package http

import "time"

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type loginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type registerRequest struct {
	Email    string `json:"email" example:"newuser@example.com"`
	Password string `json:"password" example:"secret123"`
}
