package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"orderProcessor/internal/domain"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Postgres struct {
	db PgxIface
}

func NewPostgres(db PgxIface) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) CreateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := p.createOrderMain(ctx, tx, order); err != nil {
		return err
	}

	if err := p.createOrderDelivery(ctx, tx, order.OrderUID, &order.Delivery); err != nil {
		return err
	}

	if err := p.createOrderPayment(ctx, tx, order.OrderUID, &order.Payment); err != nil {
		return err
	}

	if err := p.createOrderItems(ctx, tx, order.OrderUID, order.Items); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *Postgres) createOrderMain(ctx context.Context, tx pgx.Tx, order *domain.Order) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	return err
}

func (p *Postgres) createOrderDelivery(ctx context.Context, tx pgx.Tx, orderUID string, delivery *domain.Delivery) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email)
	return err
}

func (p *Postgres) createOrderPayment(ctx context.Context, tx pgx.Tx, orderUID string, payment *domain.Payment) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal,
		payment.CustomFee)
	return err
}

func (p *Postgres) createOrderItems(ctx context.Context, tx pgx.Tx, orderUID string, items []domain.Item) error {
	for _, item := range items {
		_, err := tx.Exec(ctx, `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) ReadOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	order := &domain.Order{}

	row := p.db.QueryRow(ctx, `
        SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders
        WHERE order_uid = $1`, orderUID)
	if err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard); err != nil {
		return nil, err
	}

	if err := p.readOrderDelivery(ctx, orderUID, order); err != nil {
		return nil, err
	}

	if err := p.readOrderPayment(ctx, orderUID, order); err != nil {
		return nil, err
	}

	if err := p.readOrderItems(ctx, orderUID, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (p *Postgres) readOrderDelivery(ctx context.Context, orderUID string, order *domain.Order) error {
	row := p.db.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1`, orderUID)
	return row.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
}

func (p *Postgres) readOrderPayment(ctx context.Context, orderUID string, order *domain.Order) error {
	row := p.db.QueryRow(ctx, `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1`, orderUID)
	return row.Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee)
}

func (p *Postgres) readOrderItems(ctx context.Context, orderUID string, order *domain.Order) error {
	rows, err := p.db.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1`, orderUID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		item := domain.Item{}
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			return err
		}
		order.Items = append(order.Items, item)
	}

	return nil
}

func (p *Postgres) ReadAllOrders(ctx context.Context) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0)

	rows, err := p.db.Query(ctx, `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := domain.Order{}
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard); err != nil {
			return nil, err
		}

		if err := p.readOrderDelivery(ctx, order.OrderUID, &order); err != nil {
			return nil, err
		}

		if err := p.readOrderPayment(ctx, order.OrderUID, &order); err != nil {
			return nil, err
		}

		if err := p.readOrderItems(ctx, order.OrderUID, &order); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}
