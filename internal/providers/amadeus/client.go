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

const (
	providerName   = "Amadeus"
	tokenEndpoint  = "/v1/security/oauth2/token"
	offersEndpoint = "/v2/shopping/flight-offers"
	dateLayout     = "2006-01-02"
	defaultTimeout = 10 * time.Second
)

type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	maxResults string
	http       *http.Client
	token      *token
}

type token struct {
	access  string
	expires time.Time
}

// New constructs an Amadeus provider client.
func New(apiKey, apiSecret, rawURL, maxResults string, httpClient *http.Client) providers.Provider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		baseURL:    strings.TrimRight(rawURL, "/"),
		maxResults: maxResults,
		http:       httpClient,
	}
}

// GetFlights retrieves flight offers from Amadeus.
func (c *Client) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, fmt.Errorf("%s: auth failed: %w", providerName, err)
	}

	// Prepare query parameters
	q := url.Values{
		"originLocationCode":      {search.Origin},
		"destinationLocationCode": {search.Destination},
		"departureDate":           {search.DepartureDate.Format(dateLayout)},
		"adults":                  {"1"},
		"nonStop":                 {"false"},
		"max":                     {c.maxResults},
	}

	// Execute request
	var resp models.AmadeusFlightResponse
	if err := c.doRequest(ctx, http.MethodGet, offersEndpoint, q, &resp); err != nil {
		return nil, err
	}

	return c.mapOffers(resp.Data), nil
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, query url.Values, out interface{}) error {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return fmt.Errorf("%s: invalid path %s: %w", providerName, endpoint, err)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%s: request build failed: %w", providerName, err)
	}
	req.Header.Set("Accept", "application/json")
	if endpoint == offersEndpoint {
		req.Header.Set("Authorization", "Bearer "+c.token.access)
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%s: call failed: %w", providerName, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%s: read failed: %w", providerName, err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: status %d: %s", providerName, res.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("%s: decode failed: %w", providerName, err)
	}
	return nil
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != nil && time.Now().Before(c.token.expires) {
		return nil
	}
	// Refresh token
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.apiKey},
		"client_secret": {c.apiSecret},
	}
	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := c.doRequest(ctx, http.MethodPost, tokenEndpoint, form, &tr); err != nil {
		return err
	}
	c.token = &token{
		access:  tr.AccessToken,
		expires: time.Now().Add(time.Duration(tr.ExpiresIn-30) * time.Second),
	}
	return nil
}

// mapOffers converts raw API data into FlightOffer models.
func (c *Client) mapOffers(data []models.AmadeusFlightOffer) []models.FlightOffer {
	offers := make([]models.FlightOffer, 0, len(data))
	for _, d := range data {
		if len(d.Itineraries) == 0 || len(d.Itineraries[0].Segments) == 0 {
			continue
		}
		seg := d.Itineraries[0].Segments[0]
		offers = append(offers, models.FlightOffer{
			Provider:    providerName,
			Price:       parsePrice(d.Price.Total),
			Duration:    d.Itineraries[0].Duration,
			Origin:      seg.Departure.IataCode,
			Destination: seg.Arrival.IataCode,
			Date:        seg.Departure.At[:10],
		})
	}
	return offers
}

// parsePrice safely converts price string to float64.
func parsePrice(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
