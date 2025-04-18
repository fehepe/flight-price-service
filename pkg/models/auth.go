package models

// TokenRequest is the expected JSON payload to generate a JWT.
type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse is returned after successful authentication.
type TokenResponse struct {
	Token string `json:"token"`
}
