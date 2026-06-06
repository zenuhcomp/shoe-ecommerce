package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

func (r *ContentRepo) Site(ctx context.Context) (models.SiteConfig, error) {
	var raw string
	err := r.db.QueryRowContext(ctx, `SELECT data FROM site_config WHERE id = 1`).Scan(&raw)
	if err != nil {
		if err == sql.ErrNoRows {
			site := defaultSiteConfig()
			if seedErr := r.UpdateSite(ctx, site); seedErr != nil {
				return models.SiteConfig{}, fmt.Errorf("seed site config: %w", seedErr)
			}
			return site, nil
		}
		return models.SiteConfig{}, fmt.Errorf("site config: %w", err)
	}
	var site models.SiteConfig
	if err := json.Unmarshal([]byte(raw), &site); err != nil {
		return models.SiteConfig{}, err
	}
	return site, nil
}

func (r *ContentRepo) UpdateSite(ctx context.Context, site models.SiteConfig) error {
	raw, err := json.Marshal(site)
	if err != nil {
		return err
	}
	ts := nowUTC()
	var q string
	if r.driver == "pgx" {
		q = `INSERT INTO site_config (id, data, updated_at) VALUES (1, $1, $2)
ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data, updated_at = EXCLUDED.updated_at`
	} else {
		q = `INSERT OR REPLACE INTO site_config (id, data, updated_at) VALUES (1, ?, ?)`
	}
	_, err = r.db.ExecContext(ctx, q, string(raw), ts)
	return err
}

func (r *ContentRepo) BackfillSitePayment(ctx context.Context) error {
	site, err := r.Site(ctx)
	if err != nil {
		return err
	}
	if site.Payment.TillNumber != "" {
		return nil
	}
	def := defaultSiteConfig()
	site.Payment = def.Payment
	return r.UpdateSite(ctx, site)
}

func (r *ContentRepo) SiteUpdatedAt(ctx context.Context) (time.Time, error) {
	var ts time.Time
	err := r.db.QueryRowContext(ctx, `SELECT updated_at FROM site_config WHERE id = 1`).Scan(&ts)
	return ts, err
}
