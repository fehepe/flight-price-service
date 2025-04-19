package cache

import (
	"context"
	"sync"

	"github.com/fehepe/flight-price-service/pkg/models"
)

type MockCache struct {
	store map[string][]models.FlightOffer
	mu    sync.RWMutex
}

func NewMockCache() *MockCache {
	return &MockCache{
		store: make(map[string][]models.FlightOffer),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) ([]models.FlightOffer, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.store[key]
	return data, ok, nil
}

func (m *MockCache) Set(ctx context.Context, key string, offers []models.FlightOffer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = offers
	return nil
}
