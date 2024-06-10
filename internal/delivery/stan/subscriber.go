package natsstr

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"orderProcessor/internal/domain"
	"orderProcessor/pkg/validate"
	"time"
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
	_, err := sc.Subscribe("order.create", s.messageHandler, stan.SetManualAckMode(), stan.AckWait(60*time.Second), stan.MaxInflight(1), stan.DeliverAllAvailable())
	return err
}

func (s *Subscriber) messageHandler(m *stan.Msg) {
	var order domain.Order
	if err := json.Unmarshal(m.Data, &order); err != nil {
		log.Println("failed to unmarshal message", err)
		return
	}
	log.Println("received order", order)

	if err := validate.ValidateStruct(order); err != nil {
		log.Println("failed to validate order", err)
		m.Ack()
		log.Println("acknowledged message")
		return
	}

	if err := s.repo.Create(context.Background(), &order); err != nil {
		log.Println("failed to create order in repository", err)
	}
	m.Ack()
}
