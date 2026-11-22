package repository

import (
	"context"
	"database/sql"
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

func (r *ContentRepo) ListOrders(ctx context.Context, limit int) ([]models.OrderRecord, error) {
	if limit <= 0 {
		limit = 100
	}
	q := rebind(r.driver, `
SELECT id, name, phone, email, service_slug, service_name, size, quantity, amount_kes,
       payment_method, mpesa_phone, offline_channel, payment_reference, status, notes, created_at
FROM orders ORDER BY created_at DESC LIMIT ?`)
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.OrderRecord
	for rows.Next() {
		var o models.OrderRecord
		var email, size, mpesa, offline, ref, notes sql.NullString
		if err := rows.Scan(
			&o.ID, &o.Name, &o.Phone, &email, &o.ServiceSlug, &o.ServiceName,
			&size, &o.Quantity, &o.AmountKES, &o.PaymentMethod,
			&mpesa, &offline, &ref, &o.Status, &notes, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		o.Email = email.String
		o.Size = size.String
		o.MpesaPhone = mpesa.String
		o.OfflineChannel = offline.String
		o.PaymentReference = ref.String
		o.Notes = notes.String
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *ContentRepo) UpdateOrderStatus(ctx context.Context, id, status string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `UPDATE orders SET status = ? WHERE id = ?`), status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *ContentRepo) DeleteOrder(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM orders WHERE id = ?`), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *ContentRepo) GetOrderByID(ctx context.Context, id string) (*models.OrderRecord, error) {
	q := rebind(r.driver, `
SELECT id, name, phone, email, service_slug, service_name, size, quantity, amount_kes,
       payment_method, mpesa_phone, offline_channel, payment_reference, status, notes, created_at
FROM orders WHERE id = ? LIMIT 1`)
	var o models.OrderRecord
	var email, size, mpesa, offline, ref, notes sql.NullString
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&o.ID, &o.Name, &o.Phone, &email, &o.ServiceSlug, &o.ServiceName,
		&size, &o.Quantity, &o.AmountKES, &o.PaymentMethod,
		&mpesa, &offline, &ref, &o.Status, &notes, &o.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", id)
		}
		return nil, err
	}
	o.Email = email.String
	o.Size = size.String
	o.MpesaPhone = mpesa.String
	o.OfflineChannel = offline.String
	o.PaymentReference = ref.String
	o.Notes = notes.String
	return &o, nil
}

func (r *ContentRepo) UpdateOrderPaymentRef(ctx context.Context, orderID, checkoutID string) error {
	q := rebind(r.driver, `UPDATE orders SET payment_reference = ? WHERE id = ?`)
	_, err := r.db.ExecContext(ctx, q, checkoutID, orderID)
	return err
}

func (r *ContentRepo) ReconcileOrderPayment(ctx context.Context, checkoutID, receipt, status, noteAppend string) (*models.OrderRecord, error) {
	qFind := rebind(r.driver, `
SELECT id, name, phone, email, service_slug, service_name, size, quantity, amount_kes,
       payment_method, mpesa_phone, offline_channel, payment_reference, status, notes, created_at
FROM orders WHERE payment_reference = ? OR id = ? ORDER BY created_at DESC LIMIT 1`)
	var o models.OrderRecord
	var email, size, mpesa, offline, ref, notes sql.NullString
	err := r.db.QueryRowContext(ctx, qFind, checkoutID, checkoutID).Scan(
		&o.ID, &o.Name, &o.Phone, &email, &o.ServiceSlug, &o.ServiceName,
		&size, &o.Quantity, &o.AmountKES, &o.PaymentMethod,
		&mpesa, &offline, &ref, &o.Status, &notes, &o.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no matching order found for checkout/ref: %s", checkoutID)
		}
		return nil, err
	}

	newNotes := notes.String
	if noteAppend != "" {
		if newNotes != "" {
			newNotes += "\n" + noteAppend
		} else {
			newNotes = noteAppend
		}
	}

	qUpdate := rebind(r.driver, `UPDATE orders SET status = ?, payment_reference = ?, notes = ? WHERE id = ?`)
	if _, err := r.db.ExecContext(ctx, qUpdate, status, receipt, newNotes, o.ID); err != nil {
		return nil, fmt.Errorf("update order %s: %w", o.ID, err)
	}

	o.Status = status
	o.PaymentReference = receipt
	o.Notes = newNotes
	return &o, nil
}

func (r *ContentRepo) nextOrderID(ctx context.Context) (string, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&count); err != nil {
		return "", err
	}
	return fmt.Sprintf("ord-%05d", count+1), nil
}
