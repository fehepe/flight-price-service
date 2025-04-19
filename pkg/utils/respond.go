package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fehepe/flight-price-service/pkg/models"
)

// RespondError writes a standardized JSON error with given status and message.
func RespondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse{Error: message})
}

// RespondJSON writes a JSON response with given status and payload.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	respondJSON(w, status, payload)
}

// respondJSON is a helper that writes JSON with proper headers and logs on failure.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}
