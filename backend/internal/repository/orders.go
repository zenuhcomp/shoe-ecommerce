package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

func (r *ContentRepo) AddOrder(ctx context.Context, order models.OrderRecord) (string, error) {
	id, err := r.nextOrderID(ctx)
	if err != nil {
		return "", err
	}
	ts := time.Now().UTC()
	q := rebind(r.driver, `
INSERT INTO orders (
	id, name, phone, email, service_slug, service_name, size, quantity, amount_kes,
	payment_method, mpesa_phone, offline_channel, payment_reference, status, notes, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if r.driver == "pgx" {
		q = `
INSERT INTO orders (
	id, name, phone, email, service_slug, service_name, size, quantity, amount_kes,
	payment_method, mpesa_phone, offline_channel, payment_reference, status, notes, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	}
	_, err = r.db.ExecContext(ctx, q,
		id, order.Name, order.Phone, nullStr(order.Email), order.ServiceSlug, order.ServiceName,
		nullStr(order.Size), order.Quantity, order.AmountKES, order.PaymentMethod,
		nullStr(order.MpesaPhone), nullStr(order.OfflineChannel), nullStr(order.PaymentReference),
		order.Status, nullStr(order.Notes), ts,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ContentRepo) nextOrderID(ctx context.Context) (string, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&count); err != nil {
		return "", err
	}
	return fmt.Sprintf("ord-%05d", count+1), nil
}
