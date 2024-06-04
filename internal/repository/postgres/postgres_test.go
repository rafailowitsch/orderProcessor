package postgres

import (
	"context"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"orderProcessor/internal/domain"
	"reflect"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	// Создаем моковые ожидания для всех запросов
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO orders").WithArgs(
		"order_uid_example", "track_number_example", "entry_example", "locale_example", "",
		"customer_id_example", "delivery_service_example", "shardkey_example", 1, "2021-11-26T06:22:19Z", "1",
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO delivery").WithArgs(
		"order_uid_example", "name_example", "phone_example", "zip_example", "city_example",
		"address_example", "region_example", "email_example",
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO payment").WithArgs(
		"order_uid_example", "transaction_example", "", "USD", "provider_example", 100,
		1637907727, "bank_example", 1500, 317, 0,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO items").WithArgs(
		"order_uid_example", 9934930, "track_number_example", 453, "rid_example", "name_example",
		30, "0", 317, 2389212, "brand_example", 202,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	// Пример данных для тестирования
	order := &domain.Order{
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

	err = repo.CreateOrder(ctx, order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrderMain(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	order := &domain.Order{
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
	}

	mock.ExpectExec("INSERT INTO orders").WithArgs("order_uid_example", "track_number_example", "entry_example", "locale_example", "",
		"customer_id_example", "delivery_service_example", "shardkey_example", 1, "2021-11-26T06:22:19Z", "1").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.createOrderMain(ctx, mock, order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrderDelivery(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	delivery := &domain.Delivery{
		Name:    "name_example",
		Phone:   "phone_example",
		Zip:     "zip_example",
		City:    "city_example",
		Address: "address_example",
		Region:  "region_example",
		Email:   "email_example",
	}

	mock.ExpectExec("INSERT INTO delivery").WithArgs("order_uid_example", "name_example", "phone_example", "zip_example",
		"city_example", "address_example", "region_example", "email_example").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.createOrderDelivery(ctx, mock, "order_uid_example", delivery)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrderPayment(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	payment := &domain.Payment{
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
	}

	mock.ExpectExec("INSERT INTO payment").WithArgs("order_uid_example", "transaction_example", "",
		"USD", "provider_example", 100, 1637907727, "bank_example", 1500, 317, 0).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.createOrderPayment(ctx, mock, "order_uid_example", payment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrderItems(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	items := []domain.Item{
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
	}

	mock.ExpectExec("INSERT INTO items").WithArgs("order_uid_example", 9934930, "track_number_example",
		453, "rid_example", "name_example", 30, "0", 317, 2389212, "brand_example", 202).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.createOrderItems(ctx, mock, "order_uid_example", items)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()
	orderUID := "order_uid_example"

	// order
	mock.ExpectQuery(`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
			AddRow(orderUID, "track_number_example", "entry_example", "locale_example", "", "customer_id_example", "delivery_service_example", "shardkey_example", 1, "2021-11-26T06:22:19Z", "1"))

	// delivery
	mock.ExpectQuery(`SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"name", "phone", "zip", "city", "address", "region", "email"}).
			AddRow("name_example", "phone_example", "zip_example", "city_example", "address_example", "region_example", "email_example"))

	// payment
	mock.ExpectQuery(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}).
			AddRow("transaction_example", "", "USD", "provider_example", 100, 1637907727, "bank_example", 1500, 317, 0))

	// items
	mock.ExpectQuery(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}).
			AddRow(9934930, "track_number_example", 453, "rid_example", "name_example", 30, "0", 317, 2389212, "brand_example", 202))

	order, err := repo.ReadOrder(ctx, orderUID)
	assert.NoError(t, err)
	assert.NotNil(t, order)

	// Проверка полей order
	expectedOrderFields := map[string]interface{}{
		"OrderUID":          "order_uid_example",
		"TrackNumber":       "track_number_example",
		"Entry":             "entry_example",
		"Locale":            "locale_example",
		"InternalSignature": "",
		"CustomerID":        "customer_id_example",
		"DeliveryService":   "delivery_service_example",
		"Shardkey":          "shardkey_example",
		"SmID":              1,
		"DateCreated":       "2021-11-26T06:22:19Z",
		"OofShard":          "1",
	}

	for field, expectedValue := range expectedOrderFields {
		assert.Equal(t, expectedValue, getFieldValue(order, field), "Field %s", field)
	}

	// Проверка полей delivery
	expectedDeliveryFields := map[string]interface{}{
		"Name":    "name_example",
		"Phone":   "phone_example",
		"Zip":     "zip_example",
		"City":    "city_example",
		"Address": "address_example",
		"Region":  "region_example",
		"Email":   "email_example",
	}

	for field, expectedValue := range expectedDeliveryFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Delivery, field), "Field %s", field)
	}

	// Проверка полей payment
	expectedPaymentFields := map[string]interface{}{
		"Transaction":  "transaction_example",
		"RequestID":    "",
		"Currency":     "USD",
		"Provider":     "provider_example",
		"Amount":       100,
		"PaymentDt":    1637907727,
		"Bank":         "bank_example",
		"DeliveryCost": 1500,
		"GoodsTotal":   317,
		"CustomFee":    0,
	}

	for field, expectedValue := range expectedPaymentFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Payment, field), "Field %s", field)
	}

	// Проверка полей items
	assert.Len(t, order.Items, 1)
	item := order.Items[0]
	expectedItemFields := map[string]interface{}{
		"ChrtID":      9934930,
		"TrackNumber": "track_number_example",
		"Price":       453,
		"Rid":         "rid_example",
		"Name":        "name_example",
		"Sale":        30,
		"Size":        "0",
		"TotalPrice":  317,
		"NmID":        2389212,
		"Brand":       "brand_example",
		"Status":      202,
	}

	for field, expectedValue := range expectedItemFields {
		assert.Equal(t, expectedValue, getFieldValue(&item, field), "Field %s", field)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadOrderDelivery(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()
	orderUID := "order_uid_example"

	mock.ExpectQuery(`SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"name", "phone", "zip", "city", "address", "region", "email"}).
			AddRow("name_example", "phone_example", "zip_example", "city_example", "address_example", "region_example", "email_example"))

	order := &domain.Order{}
	err = repo.readOrderDelivery(ctx, orderUID, order)
	assert.NoError(t, err)

	expectedDeliveryFields := map[string]interface{}{
		"Name":    "name_example",
		"Phone":   "phone_example",
		"Zip":     "zip_example",
		"City":    "city_example",
		"Address": "address_example",
		"Region":  "region_example",
		"Email":   "email_example",
	}

	for field, expectedValue := range expectedDeliveryFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Delivery, field), "Field %s", field)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadOrderPayment(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()
	orderUID := "order_uid_example"

	mock.ExpectQuery(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}).
			AddRow("transaction_example", "", "USD", "provider_example", 100, 1637907727, "bank_example", 1500, 317, 0))

	order := &domain.Order{}
	err = repo.readOrderPayment(ctx, orderUID, order)
	assert.NoError(t, err)

	expectedPaymentFields := map[string]interface{}{
		"Transaction":  "transaction_example",
		"RequestID":    "",
		"Currency":     "USD",
		"Provider":     "provider_example",
		"Amount":       100,
		"PaymentDt":    1637907727,
		"Bank":         "bank_example",
		"DeliveryCost": 1500,
		"GoodsTotal":   317,
		"CustomFee":    0,
	}

	for field, expectedValue := range expectedPaymentFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Payment, field), "Field %s", field)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadOrderItems(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()
	orderUID := "order_uid_example"

	mock.ExpectQuery(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = \$1`).
		WithArgs(orderUID).
		WillReturnRows(pgxmock.NewRows([]string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}).
			AddRow(9934930, "track_number_example", 453, "rid_example", "name_example", 30, "0", 317, 2389212, "brand_example", 202))

	order := &domain.Order{}
	err = repo.readOrderItems(ctx, orderUID, order)
	assert.NoError(t, err)
	assert.Len(t, order.Items, 1)

	expectedItemFields := map[string]interface{}{
		"ChrtID":      9934930,
		"TrackNumber": "track_number_example",
		"Price":       453,
		"Rid":         "rid_example",
		"Name":        "name_example",
		"Sale":        30,
		"Size":        "0",
		"TotalPrice":  317,
		"NmID":        2389212,
		"Brand":       "brand_example",
		"Status":      202,
	}

	for field, expectedValue := range expectedItemFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Items[0], field), "Field %s", field)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadAllOrders(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock connection: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewPostgres(mock)

	ctx := context.Background()

	// Создаем моковые ожидания для всех запросов
	mock.ExpectQuery(`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`).
		WillReturnRows(pgxmock.NewRows([]string{"order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}).
			AddRow("order_uid_example", "track_number_example", "entry_example", "locale_example", "", "customer_id_example", "delivery_service_example", "shardkey_example", 1, "2021-11-26T06:22:19Z", "1"))

	mock.ExpectQuery(`SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = \$1`).
		WithArgs("order_uid_example").
		WillReturnRows(pgxmock.NewRows([]string{"name", "phone", "zip", "city", "address", "region", "email"}).
			AddRow("name_example", "phone_example", "zip_example", "city_example", "address_example", "region_example", "email_example"))

	mock.ExpectQuery(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = \$1`).
		WithArgs("order_uid_example").
		WillReturnRows(pgxmock.NewRows([]string{"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}).
			AddRow("transaction_example", "", "USD", "provider_example", 100, 1637907727, "bank_example", 1500, 317, 0))

	mock.ExpectQuery(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = \$1`).
		WithArgs("order_uid_example").
		WillReturnRows(pgxmock.NewRows([]string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}).
			AddRow(9934930, "track_number_example", 453, "rid_example", "name_example", 30, "0", 317, 2389212, "brand_example", 202))

	orders, err := repo.ReadAllOrders(ctx)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)

	order := orders[0]

	// Проверка полей order
	expectedOrderFields := map[string]interface{}{
		"OrderUID":          "order_uid_example",
		"TrackNumber":       "track_number_example",
		"Entry":             "entry_example",
		"Locale":            "locale_example",
		"InternalSignature": "",
		"CustomerID":        "customer_id_example",
		"DeliveryService":   "delivery_service_example",
		"Shardkey":          "shardkey_example",
		"SmID":              1,
		"DateCreated":       "2021-11-26T06:22:19Z",
		"OofShard":          "1",
	}

	for field, expectedValue := range expectedOrderFields {
		assert.Equal(t, expectedValue, getFieldValue(order, field), "Field %s", field)
	}

	// Проверка полей delivery
	expectedDeliveryFields := map[string]interface{}{
		"Name":    "name_example",
		"Phone":   "phone_example",
		"Zip":     "zip_example",
		"City":    "city_example",
		"Address": "address_example",
		"Region":  "region_example",
		"Email":   "email_example",
	}

	for field, expectedValue := range expectedDeliveryFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Delivery, field), "Field %s", field)
	}

	// Проверка полей payment
	expectedPaymentFields := map[string]interface{}{
		"Transaction":  "transaction_example",
		"RequestID":    "",
		"Currency":     "USD",
		"Provider":     "provider_example",
		"Amount":       100,
		"PaymentDt":    1637907727,
		"Bank":         "bank_example",
		"DeliveryCost": 1500,
		"GoodsTotal":   317,
		"CustomFee":    0,
	}

	for field, expectedValue := range expectedPaymentFields {
		assert.Equal(t, expectedValue, getFieldValue(&order.Payment, field), "Field %s", field)
	}

	// Проверка полей items
	assert.Len(t, order.Items, 1)
	item := order.Items[0]
	expectedItemFields := map[string]interface{}{
		"ChrtID":      9934930,
		"TrackNumber": "track_number_example",
		"Price":       453,
		"Rid":         "rid_example",
		"Name":        "name_example",
		"Sale":        30,
		"Size":        "0",
		"TotalPrice":  317,
		"NmID":        2389212,
		"Brand":       "brand_example",
		"Status":      202,
	}

	for field, expectedValue := range expectedItemFields {
		assert.Equal(t, expectedValue, getFieldValue(&item, field), "Field %s", field)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func getFieldValue(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}
