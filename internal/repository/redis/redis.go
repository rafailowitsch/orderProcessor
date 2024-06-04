package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"orderProcessor/internal/domain"
)

type Redis struct {
	cache *redis.Client
}

func NewRedis(cache *redis.Client) *Redis {
	return &Redis{
		cache: cache,
	}
}

func (r *Redis) Set(ctx context.Context, order *domain.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Could not marshal JSON: %v", err)
		return err
	}

	key := "order:" + order.OrderUID

	err = r.cache.Set(ctx, key, orderJSON, 0).Err()
	if err != nil {
		log.Fatalf("Could not set JSON in Redis: %v", err)
		return err
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	key := "order:" + orderUID

	orderJSON, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		log.Fatalf("Could not get JSON from Redis: %v", err)
		return nil, err
	}

	order := &domain.Order{}
	err = json.Unmarshal([]byte(orderJSON), order)
	if err != nil {
		log.Fatalf("Could not unmarshal JSON: %v", err)
		return nil, err
	}

	return order, nil
}
