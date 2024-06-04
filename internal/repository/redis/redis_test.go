package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"orderProcessor/internal/domain"
	"testing"
)

var testOrder = &domain.Order{
	OrderUID:          "order_uid_example",
	TrackNumber:       "track_number_example",
	Entry:             "entry_example",
	Locale:            "locale_example",
	InternalSignature: "",
	CustomerID:        "customer_id_example",
	DeliveryService:   "delivery_service_example",
	Shardkey:          "shardkey_example",
	SmID:              1,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
	Delivery: domain.Delivery{
		Name:    "name_example",
		Phone:   "phone_example",
		Zip:     "zip_example",
		City:    "city_example",
		Address: "address_example",
		Region:  "region_example",
		Email:   "email_example",
	},
	Payment: domain.Payment{
		Transaction:  "transaction_example",
		RequestID:    "",
		Currency:     "USD",
		Provider:     "provider_example",
		Amount:       100,
		PaymentDt:    1637907727,
		Bank:         "bank_example",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []domain.Item{
		{
			ChrtID:      9934930,
			TrackNumber: "track_number_example",
			Price:       453,
			Rid:         "rid_example",
			Name:        "name_example",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "brand_example",
			Status:      202,
		},
	},
}

func TestRedis_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	r := NewRedis(db)

	ctx := context.Background()
	order := testOrder
	orderJSON, err := json.Marshal(order)
	assert.NoError(t, err)

	mock.ExpectSet("order:order_uid_example", orderJSON, 0).SetVal("OK")

	err = r.Set(ctx, order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedis_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	r := NewRedis(db)

	ctx := context.Background()
	orderUID := "order_uid_example"
	expectedOrder := testOrder

	orderJSON, err := json.Marshal(expectedOrder)
	assert.NoError(t, err)

	mock.ExpectGet("order:" + orderUID).SetVal(string(orderJSON))

	order, err := r.Get(ctx, orderUID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	assert.NoError(t, mock.ExpectationsWereMet())
}
