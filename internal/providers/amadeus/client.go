package amadeus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

type token struct {
	AccessToken string
	ExpiresAt   time.Time
}

type Client struct {
	apiKey           string
	apiSecret        string
	baseURL          string
	maxFlightResults string
	httpClient       *http.Client
	token            *token
}

func New(apiKey, apiSecret, baseURL, maxResults string, httpClient *http.Client) providers.Provider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		apiKey:           apiKey,
		apiSecret:        apiSecret,
		baseURL:          baseURL,
		maxFlightResults: maxResults,
		httpClient:       httpClient,
	}
}

func (c *Client) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("token retrieval failed: %w", err)
	}

	u, err := url.Parse(c.baseURL + "/v2/shopping/flight-offers")
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	params := url.Values{}
	params.Set("originLocationCode", search.Origin)
	params.Set("destinationLocationCode", search.Destination)
	params.Set("departureDate", search.DepartureDate.Format("2006-01-02"))
	params.Set("adults", "1")
	params.Set("max", c.maxFlightResults)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("amadeus API error [%d]: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var result models.AmadeusFlightResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("response decode failed: %w", err)
	}

	offers := make([]models.FlightOffer, 0, len(result.Data))
	for _, d := range result.Data {
		if len(d.Itineraries) == 0 || len(d.Itineraries[0].Segments) == 0 {
			continue
		}
		seg := d.Itineraries[0].Segments[0]
		offers = append(offers, models.FlightOffer{
			Provider:    "Amadeus",
			Price:       parsePrice(d.Price.Total),
			Duration:    d.Itineraries[0].Duration,
			Origin:      seg.Departure.IataCode,
			Destination: seg.Arrival.IataCode,
			Date:        seg.Departure.At[:10],
		})
	}

	return offers, nil
}

func (c *Client) getToken() (string, error) {
	if c.token != nil && time.Now().Before(c.token.ExpiresAt) {
		return c.token.AccessToken, nil
	}
	return c.fetchNewToken()
}

func (c *Client) fetchNewToken() (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.apiKey)
	form.Set("client_secret", c.apiSecret)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/security/oauth2/token", c.baseURL), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("token error [%d]: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("token decode failed: %w", err)
	}

	t := &token{
		AccessToken: tr.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(tr.ExpiresIn-30) * time.Second),
	}
	c.token = t
	return t.AccessToken, nil
}

func parsePrice(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
