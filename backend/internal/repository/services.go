package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

const serviceSelectCols = `id, name, slug, emoji, style, length, description, price_kes, featured, images, view_360, category, badge, original_price_kes, currency, specs, video_url, in_stock`

func scanService(scanner interface {
	Scan(dest ...any) error
}) (models.LashService, error) {
	var s models.LashService
	var featured, inStock int
	var imagesRaw, view360Raw, specsRaw string
	var origPriceKES int
	var category, badge, currency, videoURL string
	if err := scanner.Scan(
		&s.ID, &s.Name, &s.Slug, &s.Emoji, &s.Style, &s.Length, &s.Description, &s.PriceKES, &featured,
		&imagesRaw, &view360Raw,
		&category, &badge, &origPriceKES, &currency, &specsRaw, &videoURL, &inStock,
	); err != nil {
		return s, err
	}
	s.Featured = featured == 1
	s.InStock = inStock == 1
	s.Category = category
	if s.Category == "" && s.Style != "" {
		s.Category = s.Style
	}
	s.Badge = badge
	s.OriginalPriceKES = origPriceKES
	s.Currency = currency
	if s.Currency == "" {
		s.Currency = "KES"
	}
	s.VideoURL = videoURL
	s.Images = unmarshalJSONStrings(imagesRaw)
	s.View360 = unmarshalJSONStrings(view360Raw)
	if len(s.View360) == 0 {
		s.View360 = append([]string{}, s.Images...)
	}
	s.Specs = unmarshalJSONMap(specsRaw)
	return s, nil
}

func (r *ContentRepo) Services(ctx context.Context) ([]models.LashService, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s FROM lash_services ORDER BY sort_order ASC, name ASC`, serviceSelectCols))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.LashService
	for rows.Next() {
		s, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *ContentRepo) ServiceBySlug(ctx context.Context, slug string) (models.LashService, bool, error) {
	rows, err := r.Services(ctx)
	if err != nil {
		return models.LashService{}, false, err
	}
	for _, s := range rows {
		if s.Slug == slug {
			return s, true, nil
		}
	}
	return models.LashService{}, false, nil
}

func (r *ContentRepo) ServiceByID(ctx context.Context, id string) (models.LashService, error) {
	q := rebind(r.driver, fmt.Sprintf(`
SELECT %s FROM lash_services WHERE id = ?`, serviceSelectCols))
	s, err := scanService(r.db.QueryRowContext(ctx, q, id))
	if err == sql.ErrNoRows {
		return s, fmt.Errorf("not found")
	}
	if err != nil {
		return s, err
	}
	return s, nil
}

func (r *ContentRepo) CreateService(ctx context.Context, s models.LashService, sortOrder int) error {
	ts := nowUTC()
	if s.Currency == "" {
		s.Currency = "KES"
	}
	inStock := 1
	if !s.InStock && s.Category != "" {
		inStock = 0
	}
	q := rebind(r.driver, `
INSERT INTO lash_services (
	id, name, slug, emoji, style, length, description, price_kes, featured, images, view_360,
	category, badge, original_price_kes, currency, specs, video_url, in_stock, sort_order, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := r.db.ExecContext(ctx, q,
		s.ID, s.Name, s.Slug, s.Emoji, s.Style, s.Length, s.Description, s.PriceKES, boolInt(s.Featured),
		marshalJSONStrings(s.Images), marshalJSONStrings(s.View360),
		s.Category, s.Badge, s.OriginalPriceKES, s.Currency, marshalJSONMap(s.Specs), s.VideoURL, inStock,
		sortOrder, ts, ts,
	)
	return err
}

func (r *ContentRepo) UpdateService(ctx context.Context, s models.LashService) error {
	ts := nowUTC()
	if s.Currency == "" {
		s.Currency = "KES"
	}
	q := rebind(r.driver, `
UPDATE lash_services SET
	name=?, slug=?, emoji=?, style=?, length=?, description=?, price_kes=?, featured=?, images=?, view_360=?,
	category=?, badge=?, original_price_kes=?, currency=?, specs=?, video_url=?, in_stock=?, updated_at=?
WHERE id=?`)
	res, err := r.db.ExecContext(ctx, q,
		s.Name, s.Slug, s.Emoji, s.Style, s.Length, s.Description, s.PriceKES, boolInt(s.Featured),
		marshalJSONStrings(s.Images), marshalJSONStrings(s.View360),
		s.Category, s.Badge, s.OriginalPriceKES, s.Currency, marshalJSONMap(s.Specs), s.VideoURL, boolInt(s.InStock), ts,
		s.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeleteService(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM lash_services WHERE id = ?`), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) ReorderServices(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "lash_services", ids)
}

func (r *ContentRepo) DeleteAllServices(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lash_services`)
	return err
}

func (r *ContentRepo) reorderStrings(ctx context.Context, table string, ids []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, id := range ids {
		q := rebind(r.driver, fmt.Sprintf(`UPDATE %s SET sort_order = ? WHERE id = ?`, table))
		if _, err := tx.ExecContext(ctx, q, i, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

var staleSeedImageFragments = []string{
	"photo-1560769629-975ec794e0a6",
	"photo-1515346819254-7a4533097a48",
	"photo-1525966222134-fcfa99b8d057",
	"photo-1533867610662-5c741a7365a2",
	"photo-1582897085656-c636d0063575",
	"photo-1460353581641-37baddab0fa2",
	"photo-1572804013309-59a881b4e889",
	"photo-1605812860427-4024433a70fd",
}

func (r *ContentRepo) BackfillServiceMedia(ctx context.Context) error {
	seedByID := map[string]models.LashService{}
	for _, s := range defaultServices() {
		seedByID[s.ID] = s
	}
	existing, err := r.Services(ctx)
	if err != nil {
		return err
	}
	for _, s := range existing {
		if len(s.Images) > 0 && !containsStaleSeedMedia(s.Images) && !containsStaleSeedMedia(s.View360) {
			continue
		}
		seed, ok := seedByID[s.ID]
		if !ok {
			continue
		}
		s.Images = seed.Images
		s.View360 = seed.View360
		if err := r.UpdateService(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func containsStaleSeedMedia(urls []string) bool {
	for _, url := range urls {
		for _, fragment := range staleSeedImageFragments {
			if strings.Contains(url, fragment) {
				return true
			}
		}
	}
	return false
}

func (r *ContentRepo) NextServiceSort(ctx context.Context) (int, error) {
	var max sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT MAX(sort_order) FROM lash_services`).Scan(&max)
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 0, nil
	}
	return int(max.Int64) + 1, nil
}
