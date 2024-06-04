package repository

import (
	"context"
	"orderProcessor/internal/domain"
)

type Database interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	ReadOrder(ctx context.Context, orderUID string) (*domain.Order, error)
	ReadAllOrders(ctx context.Context) ([]*domain.Order, error)
}

type Cache interface {
	Set(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, orderUID string) (*domain.Order, error)
}

type Repository struct {
	db    Database
	cache Cache
}

func NewRepository(db Database, cache Cache) *Repository {
	return &Repository{
		db:    db,
		cache: cache,
	}
}
