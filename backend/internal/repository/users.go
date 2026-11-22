package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kibsoft/luna-collections/backend/internal/models"
)

type cmsUserRow struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *ContentRepo) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cms_users`).Scan(&n)
	return n, err
}

func (r *ContentRepo) CreateUser(ctx context.Context, email, passwordHash, displayName string) (models.AdminUser, error) {
	id := uuid.NewString()
	ts := nowUTC()
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO cms_users (id, email, password_hash, display_name, role, created_at, updated_at)
VALUES (?, ?, ?, ?, 'admin', ?, ?)`),
		id, email, passwordHash, displayName, ts, ts)
	if err != nil {
		return models.AdminUser{}, err
	}
	return models.AdminUser{ID: id, Email: email, DisplayName: displayName, Role: "admin", CreatedAt: ts, UpdatedAt: ts}, nil
}

func (r *ContentRepo) UpsertAdminUser(ctx context.Context, email, passwordHash, displayName string) (models.AdminUser, error) {
	ts := nowUTC()
	var existingID string
	err := r.db.QueryRowContext(ctx, rebind(r.driver, `SELECT id FROM cms_users WHERE LOWER(email) = LOWER(?)`), email).Scan(&existingID)
	if err == nil {
		_, err = r.db.ExecContext(ctx, rebind(r.driver, `
UPDATE cms_users SET password_hash = ?, display_name = ?, updated_at = ? WHERE id = ?`),
			passwordHash, displayName, ts, existingID)
		if err != nil {
			return models.AdminUser{}, err
		}
		return models.AdminUser{ID: existingID, Email: email, DisplayName: displayName, Role: "admin", UpdatedAt: ts}, nil
	}

	id := uuid.NewString()
	_, err = r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO cms_users (id, email, password_hash, display_name, role, created_at, updated_at)
VALUES (?, ?, ?, ?, 'admin', ?, ?)`),
		id, email, passwordHash, displayName, ts, ts)
	if err != nil {
		return models.AdminUser{}, err
	}
	return models.AdminUser{ID: id, Email: email, DisplayName: displayName, Role: "admin", CreatedAt: ts, UpdatedAt: ts}, nil
}

func (r *ContentRepo) UserByEmail(ctx context.Context, email string) (cmsUserRow, error) {
	var u cmsUserRow
	q := rebind(r.driver, `
SELECT id, email, password_hash, display_name, role, created_at, updated_at FROM cms_users WHERE LOWER(email) = LOWER(?)`)
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return u, fmt.Errorf("not found")
	}
	return u, err
}

func (r *ContentRepo) UserByID(ctx context.Context, id string) (models.AdminUser, error) {
	var u models.AdminUser
	q := rebind(r.driver, `
SELECT id, email, display_name, role, created_at, updated_at FROM cms_users WHERE id = ?`)
	err := r.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email, &u.DisplayName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return u, fmt.Errorf("not found")
	}
	return u, err
}
