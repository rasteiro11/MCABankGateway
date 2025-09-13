package auth

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
