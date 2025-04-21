package serpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/fehepe/flight-price-service/pkg/models"
)

const (
	engine          = "google_flights"
	defaultCurrency = "USD"
	defaultLocale   = "en"
	providerName    = "SerpAPI"
	timeLayout      = "2006-01-02 15:04"
	dateLayout      = "2006-01-02"
	defaultTimeout  = 10 * time.Second
)

var ErrNoFlights = errors.New("no flight offers found")

type SerpAPIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func New(apiKey, baseURL string, httpClient *http.Client) providers.Provider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &SerpAPIClient{
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  httpClient,
	}
}

func (c *SerpAPIClient) GetFlights(ctx context.Context, search models.FlightSearch) ([]models.FlightOffer, error) {
	respData, err := c.doSearch(ctx, search)
	if err != nil {
		return nil, err
	}
	if len(respData.BestFlights) == 0 {
		return nil, ErrNoFlights
	}
	return c.mapToOffers(respData), nil
}

func (c *SerpAPIClient) doSearch(ctx context.Context, search models.FlightSearch) (*models.SerAPIResponse, error) {
	u, err := url.Parse(c.baseURL + "/search.json")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", c.baseURL, err)
	}

	qp := u.Query()
	qp.Set("engine", engine)
	qp.Set("currency", defaultCurrency)
	qp.Set("hl", defaultLocale)
	qp.Set("api_key", c.apiKey)
	qp.Set("departure_id", search.Origin)
	qp.Set("arrival_id", search.Destination)
	qp.Set("outbound_date", search.DepartureDate.Format(dateLayout))
	qp.Set("return_date", search.DepartureDate.AddDate(0, 0, 1).Format(dateLayout))
	u.RawQuery = qp.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("serpapi error: status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var result models.SerAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

func (c *SerpAPIClient) mapToOffers(data *models.SerAPIResponse) []models.FlightOffer {
	offers := make([]models.FlightOffer, 0, len(data.BestFlights))
	for _, fg := range data.BestFlights {
		offer, err := mapOffer(fg)
		if err != nil {
			log.Printf("failed to map flight option: %v", err)
			continue
		}
		offers = append(offers, offer)
	}
	return offers
}

func mapOffer(fg models.FlightOption) (models.FlightOffer, error) {
	if len(fg.Flights) == 0 {
		return models.FlightOffer{}, fmt.Errorf("no flight segments found")
	}
	first := fg.Flights[0]
	last := fg.Flights[len(fg.Flights)-1]

	date, err := extractDate(first.DepartureAirport.Time)
	if err != nil {
		return models.FlightOffer{}, fmt.Errorf("parsing date %q: %w", first.DepartureAirport.Time, err)
	}

	return models.FlightOffer{
		Provider:    providerName,
		Price:       float64(fg.Price),
		Duration:    formatISODuration(fg.TotalDuration),
		Origin:      first.DepartureAirport.ID,
		Destination: last.ArrivalAirport.ID,
		Date:        date,
	}, nil
}

func formatISODuration(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("PT%dH%dM", hours, minutes)
}

func extractDate(ts string) (string, error) {
	t, err := time.Parse(timeLayout, ts)
	if err != nil {
		return "", err
	}
	return t.Format(dateLayout), nil
}
