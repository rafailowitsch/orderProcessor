package natsstr

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"orderProcessor/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
}

type Subscriber struct {
	repo Repository
}

func NewSubscriber(repo Repository) *Subscriber {
	return &Subscriber{repo: repo}
}

func (s *Subscriber) Subscribe(sc stan.Conn) error {
	_, err := sc.Subscribe("order.create", s.messageHandler, stan.DeliverAllAvailable())
	return err
}

func (s *Subscriber) messageHandler(m *stan.Msg) {
	var order domain.Order
	if err := json.Unmarshal(m.Data, &order); err != nil {
		return
	}

	s.repo.Create(context.Background(), &order)
}
