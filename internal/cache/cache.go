package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fehepe/flight-price-service/internal/config"
	"github.com/fehepe/flight-price-service/pkg/models"
	"github.com/redis/go-redis/v9"
)

type FlightCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

type FlightCacher interface {
	Get(ctx context.Context, key string) ([]models.FlightOffer, bool, error)
	Set(ctx context.Context, key string, offers []models.FlightOffer) error
}

func NewFlightCacheFromConfig() FlightCacher {
	return NewFlightCache(
		config.Get("REDIS_HOST", "localhost")+":"+config.Get("REDIS_PORT", "6379"),
		config.Get("REDIS_PASSWORD", ""),
		config.GetEnvInt("REDIS_DB", 0),
		time.Duration(config.GetEnvInt("REDIS_DEFAULT_TTL", 30))*time.Second,
	)
}

func NewFlightCache(addr, password string, db int, defaultTTL time.Duration) *FlightCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &FlightCache{
		client:     client,
		defaultTTL: defaultTTL,
	}
}

func (c *FlightCache) Get(ctx context.Context, key string) ([]models.FlightOffer, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("cache get error: %w", err)
	}

	var offers []models.FlightOffer
	if err := json.Unmarshal([]byte(val), &offers); err != nil {
		return nil, false, fmt.Errorf("unmarshal error: %w", err)
	}
	return offers, true, nil
}

func (c *FlightCache) Set(ctx context.Context, key string, offers []models.FlightOffer) error {
	data, err := json.Marshal(offers)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	return c.client.Set(ctx, key, data, c.defaultTTL).Err()
}
