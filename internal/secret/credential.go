package secret

import (
	"encoding/json"
	"fmt"
	"os"
)

type Creds struct {
	AmadeusAPIKey    string `json:"AMADEUS_API_KEY"`
	AmadeusAPISecret string `json:"AMADEUS_API_SECRET"`
	SerAPIKey        string `json:"SER_API_KEY"`
}

func LoadCreds(path string) (*Creds, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read creds: %w", err)
	}
	var c Creds
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal creds: %w", err)
	}
	return &c, nil
}
