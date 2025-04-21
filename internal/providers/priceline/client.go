package priceline

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
	providerName   = "PriceLine"
	dateLayout     = "2006-01-02"
	defaultTimeout = 10 * time.Second
)

// Client wraps the PriceLine API.
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// New creates a new PriceLine client. Returns error if baseURL is invalid.
func New(apiKey, baseURL string, httpClient *http.Client) providers.Provider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  httpClient,
	}
}

// ErrNoFlights is returned when no flight offers are found.
var ErrNoFlights = fmt.Errorf("no flight offers found")

func (c *Client) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	listings, err := c.fetchListings(ctx, search)
	if err != nil {
		return nil, err
	}
	if len(listings) == 0 {
		return nil, ErrNoFlights
	}
	return c.mapToOffers(listings), nil
}

func (c *Client) fetchListings(ctx context.Context, search models.FlightSearch) ([]models.Listing, error) {
	u, err := url.Parse(c.baseURL + "/flights/search-one-way")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", c.baseURL, err)
	}
	q := u.Query()
	q.Set("originAirportCode", search.Origin)
	q.Set("destinationAirportCode", search.Destination)
	q.Set("departureDate", search.DepartureDate.Format(dateLayout))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("x-rapidapi-host", u.Host)
	req.Header.Set("x-rapidapi-key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var apiResp models.PriceLineAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return apiResp.Data.Listings, nil
}

// mapToOffers converts API listings into our FlightOffer type.
func (c *Client) mapToOffers(listings []models.Listing) []models.FlightOffer {
	offers := make([]models.FlightOffer, 0, len(listings))
	for _, l := range listings {
		if len(l.Slices) == 0 || len(l.Slices[0].Segments) == 0 || len(l.Airlines) == 0 {
			continue
		}
		seg := l.Slices[0].Segments[0]
		offers = append(offers, models.FlightOffer{
			Provider:    providerName,
			Price:       l.TotalPriceWithDecimal.Price,
			Duration:    toISO8601(l.Slices[0].DurationInMinutes),
			Origin:      seg.DepartInfo.Airport.Code,
			Destination: seg.ArrivalInfo.Airport.Code,
			Date:        seg.DepartInfo.Time.DateTime[:10],
		})
	}
	return offers
}

func toISO8601(minutesStr string) string {
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return "PT0H0M"
	}
	h := minutes / 60
	m := minutes % 60
	return fmt.Sprintf("PT%dH%dM", h, m)
}
