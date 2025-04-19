package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fehepe/flight-price-service/internal/handlers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

func setupAuthEnv() {
	os.Setenv("AUTH_USERNAME", "user")
	os.Setenv("AUTH_PASSWORD", "pass")
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_EXPIRY_HOURS", "1")
	os.Setenv("JWT_ISSUER", "test-service")
}

func TestGenerateToken_ValidCredentials(t *testing.T) {
	setupAuthEnv()

	payload := models.TokenRequest{
		Username: "user",
		Password: "pass",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.GenerateToken(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}

	var resp models.TokenResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal("expected a valid JSON response")
	}

	if resp.Token == "" {
		t.Error("expected a JWT token, got empty")
	}
}

func TestGenerateToken_InvalidCredentials(t *testing.T) {
	setupAuthEnv()

	payload := models.TokenRequest{
		Username: "invalid",
		Password: "invalid",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.GenerateToken(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rr.Code)
	}
}
