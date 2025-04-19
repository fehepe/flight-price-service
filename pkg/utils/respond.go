package utils

import (
	"encoding/json"
	"net/http"

	"github.com/fehepe/flight-price-service/pkg/models"
)

// RespondError writes a standardized JSON error with given status and message.
func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

// RespondJSON writes a JSON response with given status and payload.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
