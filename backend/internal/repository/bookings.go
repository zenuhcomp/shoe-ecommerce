package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

func (r *ContentRepo) AddBooking(ctx context.Context, inquiry models.BookingInquiry) (string, error) {
	id, err := r.nextBookingID(ctx)
	if err != nil {
		return "", err
	}
	ts := time.Now().UTC()
	q := rebind(r.driver, `
INSERT INTO booking_inquiries (id, name, phone, email, service_slug, preferred_at, message, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if r.driver == "pgx" {
		q = `
INSERT INTO booking_inquiries (id, name, phone, email, service_slug, preferred_at, message, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	}
	_, err = r.db.ExecContext(ctx, q,
		id, inquiry.Name, inquiry.Phone,
		nullStr(inquiry.Email), nullStr(inquiry.ServiceSlug),
		nullStr(inquiry.PreferredAt), nullStr(inquiry.Message), ts,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ContentRepo) ListBookings(ctx context.Context, limit int) ([]models.BookingRecord, error) {
	if limit <= 0 {
		limit = 100
	}
	q := rebind(r.driver, `
SELECT id, name, phone, email, service_slug, preferred_at, message, created_at
FROM booking_inquiries ORDER BY created_at DESC LIMIT ?`)
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.BookingRecord
	for rows.Next() {
		var b models.BookingRecord
		var email, slug, pref, msg sql.NullString
		if err := rows.Scan(&b.ID, &b.Name, &b.Phone, &email, &slug, &pref, &msg, &b.CreatedAt); err != nil {
			return nil, err
		}
		b.Email = email.String
		b.ServiceSlug = slug.String
		b.PreferredAt = pref.String
		b.Message = msg.String
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *ContentRepo) DeleteBooking(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM booking_inquiries WHERE id = ?`), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) nextBookingID(ctx context.Context) (string, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM booking_inquiries`).Scan(&count); err != nil {
		return "", err
	}
	return fmt.Sprintf("inq-%04d", count+1), nil
}

func nullStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
