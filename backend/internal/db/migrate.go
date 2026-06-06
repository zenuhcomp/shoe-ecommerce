package db

import (
	"context"
	"fmt"
	"strings"
)

func (d *DB) migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS booking_inquiries (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT,
			service_slug TEXT,
			preferred_at TEXT,
			message TEXT,
			created_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS cms_users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			display_name TEXT NOT NULL DEFAULT '',
			role TEXT NOT NULL DEFAULT 'admin',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS site_config (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			data TEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS lash_services (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			emoji TEXT NOT NULL DEFAULT '',
			style TEXT NOT NULL DEFAULT '',
			length TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			price_kes INTEGER NOT NULL DEFAULT 0,
			featured INTEGER NOT NULL DEFAULT 0,
			images TEXT NOT NULL DEFAULT '[]',
			view_360 TEXT NOT NULL DEFAULT '[]',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS features (
			id TEXT PRIMARY KEY,
			icon TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS press_mentions (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			handle TEXT NOT NULL DEFAULT '',
			initials TEXT NOT NULL DEFAULT '',
			rating INTEGER NOT NULL DEFAULT 5,
			quote TEXT NOT NULL DEFAULT '',
			published INTEGER NOT NULL DEFAULT 1,
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS faq_items (
			id INTEGER PRIMARY KEY,
			question TEXT NOT NULL,
			answer TEXT NOT NULL DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS gallery_items (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT NOT NULL DEFAULT '',
			image_url TEXT NOT NULL,
			alt TEXT NOT NULL DEFAULT '',
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS nav_links (
			id TEXT PRIMARY KEY,
			label TEXT NOT NULL,
			path TEXT NOT NULL,
			sort_order INTEGER NOT NULL DEFAULT 0,
			visible INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
	}

	for _, stmt := range statements {
		if _, err := d.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	alters := []string{
		`ALTER TABLE lash_services ADD COLUMN images TEXT NOT NULL DEFAULT '[]'`,
		`ALTER TABLE lash_services ADD COLUMN view_360 TEXT NOT NULL DEFAULT '[]'`,
		`CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT,
			service_slug TEXT NOT NULL,
			service_name TEXT NOT NULL,
			size TEXT,
			quantity INTEGER NOT NULL DEFAULT 1,
			amount_kes INTEGER NOT NULL DEFAULT 0,
			payment_method TEXT NOT NULL,
			mpesa_phone TEXT,
			offline_channel TEXT,
			payment_reference TEXT,
			status TEXT NOT NULL,
			notes TEXT,
			created_at TIMESTAMP NOT NULL
		)`,
	}
	for _, stmt := range alters {
		if _, err := d.ExecContext(ctx, stmt); err != nil {
			if !isIgnorableMigrateErr(err) {
				return fmt.Errorf("migrate alter: %w", err)
			}
		}
	}
	return nil
}

func isIgnorableMigrateErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return containsAny(msg, "duplicate column", "already exists", "Duplicate column")
}

func containsAny(s string, parts ...string) bool {
	for _, p := range parts {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
