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

func (r *Repository) Create(ctx context.Context, order *domain.Order) error {
	err := r.db.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	err = r.cache.Set(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Read(ctx context.Context, orderUID string) (*domain.Order, error) {
	order, err := r.cache.Get(ctx, orderUID)
	if err != nil {
		order, err = r.db.ReadOrder(ctx, orderUID)
		if err != nil {
			return nil, err
		}

		err = r.cache.Set(ctx, order)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return order, nil
}

func (r *Repository) ReadAll(ctx context.Context) ([]*domain.Order, error) {
	orders, err := r.db.ReadAllOrders(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) CacheRecovery(ctx context.Context) error {
	orders, err := r.db.ReadAllOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		err = r.cache.Set(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
