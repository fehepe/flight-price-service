// File: internal/handlers/auth.go
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fehepe/flight-price-service/internal/config"
	"github.com/fehepe/flight-price-service/pkg/models"
	"github.com/fehepe/flight-price-service/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken validates credentials and issues a JWT.
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	// Ensure JSON content
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// Read and decode body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req models.TokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user (replace with real validation)
	validUser := config.Get("AUTH_USERNAME", "user")
	validPass := config.Get("AUTH_PASSWORD", "pass")
	if req.Username != validUser || req.Password != validPass {
		utils.RespondError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Load JWT secret
	secret := config.Get("JWT_SECRET", "")
	if secret == "" {
		utils.RespondError(w, http.StatusInternalServerError, "server configuration error: missing JWT_SECRET")
		return
	}

	// Parse expiry
	expiresIn, err := strconv.Atoi(config.Get("JWT_EXPIRY_HOURS", "1"))
	if err != nil || expiresIn <= 0 {
		expiresIn = 1
	}

	// Create claims
	claims := jwt.RegisteredClaims{
		Subject:   req.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Hour)),
		Issuer:    config.Get("JWT_ISSUER", "flight-service"),
	}

	// Sign token
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "could not sign token")
		return
	}

	// Respond with token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	utils.RespondJSON(w, http.StatusOK, models.TokenResponse{Token: signedToken})
}
