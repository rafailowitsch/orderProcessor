package repository_test

import (
	"context"
	"errors"
	"orderProcessor/internal/domain"
	"orderProcessor/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) CreateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockDatabase) ReadOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	args := m.Called(ctx, orderUID)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockDatabase) ReadAllOrders(ctx context.Context) ([]*domain.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Order), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, orderUID string) (*domain.Order, error) {
	args := m.Called(ctx, orderUID)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func TestCreate(t *testing.T) {
	db := new(MockDatabase)
	cache := new(MockCache)
	repo := repository.NewRepository(db, cache)

	order := &domain.Order{}

	db.On("CreateOrder", mock.Anything, order).Return(nil)
	cache.On("Set", mock.Anything, order).Return(nil)

	err := repo.Create(context.Background(), order)

	assert.NoError(t, err)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestRead(t *testing.T) {
	db := new(MockDatabase)
	cache := new(MockCache)
	repo := repository.NewRepository(db, cache)

	order := &domain.Order{}
	orderUID := "orderUID"

	cache.On("Get", mock.Anything, orderUID).Return(nil, errors.New("not found"))
	db.On("ReadOrder", mock.Anything, orderUID).Return(order, nil)
	cache.On("Set", mock.Anything, order).Return(nil)

	result, err := repo.Read(context.Background(), orderUID)

	assert.NoError(t, err)
	assert.Equal(t, order, result)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestReadAll(t *testing.T) {
	db := new(MockDatabase)
	cache := new(MockCache)
	repo := repository.NewRepository(db, cache)

	orders := []*domain.Order{{}}

	db.On("ReadAllOrders", mock.Anything).Return(orders, nil)

	result, err := repo.ReadAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, orders, result)
	db.AssertExpectations(t)
	cache.AssertExpectations(t)
}
