package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

// --- Features ---

func (r *ContentRepo) Features(ctx context.Context) ([]models.Feature, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, icon, title, description FROM features ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Feature
	for rows.Next() {
		var f models.Feature
		if err := rows.Scan(&f.ID, &f.Icon, &f.Title, &f.Description); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreateFeature(ctx context.Context, f models.Feature, sortOrder int) error {
	ts := nowUTC()
	q := rebind(r.driver, `INSERT INTO features (id, icon, title, description, sort_order, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`)
	_, err := r.db.ExecContext(ctx, q, f.ID, f.Icon, f.Title, f.Description, sortOrder, ts, ts)
	return err
}

func (r *ContentRepo) UpdateFeature(ctx context.Context, f models.Feature) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `UPDATE features SET icon=?, title=?, description=?, updated_at=? WHERE id=?`),
		f.Icon, f.Title, f.Description, ts, f.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeleteFeature(ctx context.Context, id string) error {
	return r.deleteByID(ctx, "features", id)
}

func (r *ContentRepo) ReorderFeatures(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "features", ids)
}

// --- Press ---

func (r *ContentRepo) Press(ctx context.Context) ([]models.PressMention, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, note FROM press_mentions ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.PressMention
	for rows.Next() {
		var p models.PressMention
		if err := rows.Scan(&p.ID, &p.Name, &p.Note); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreatePress(ctx context.Context, p models.PressMention, sortOrder int) error {
	ts := nowUTC()
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `INSERT INTO press_mentions (id, name, note, sort_order, created_at, updated_at) VALUES (?,?,?,?,?,?)`),
		p.ID, p.Name, p.Note, sortOrder, ts, ts)
	return err
}

func (r *ContentRepo) UpdatePress(ctx context.Context, p models.PressMention) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `UPDATE press_mentions SET name=?, note=?, updated_at=? WHERE id=?`),
		p.Name, p.Note, ts, p.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeletePress(ctx context.Context, id string) error {
	return r.deleteByID(ctx, "press_mentions", id)
}

func (r *ContentRepo) ReorderPress(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "press_mentions", ids)
}

// --- Reviews ---

func (r *ContentRepo) Reviews(ctx context.Context, publicOnly bool) ([]models.Review, error) {
	q := `SELECT id, name, handle, initials, rating, quote, published FROM reviews`
	if publicOnly {
		q += ` WHERE published = 1`
	}
	q += ` ORDER BY sort_order`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Review
	for rows.Next() {
		var rev models.Review
		var pub int
		if err := rows.Scan(&rev.ID, &rev.Name, &rev.Handle, &rev.Initials, &rev.Rating, &rev.Quote, &pub); err != nil {
			return nil, err
		}
		rev.Published = pub == 1
		out = append(out, rev)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreateReview(ctx context.Context, rev models.Review, sortOrder int) error {
	ts := nowUTC()
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO reviews (id, name, handle, initials, rating, quote, published, sort_order, created_at, updated_at)
VALUES (?,?,?,?,?,?,1,?,?,?)`),
		rev.ID, rev.Name, rev.Handle, rev.Initials, rev.Rating, rev.Quote, sortOrder, ts, ts)
	return err
}

func (r *ContentRepo) UpdateReview(ctx context.Context, rev models.Review, published bool) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `
UPDATE reviews SET name=?, handle=?, initials=?, rating=?, quote=?, published=?, updated_at=? WHERE id=?`),
		rev.Name, rev.Handle, rev.Initials, rev.Rating, rev.Quote, boolInt(published), ts, rev.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeleteReview(ctx context.Context, id string) error {
	return r.deleteByID(ctx, "reviews", id)
}

func (r *ContentRepo) ReorderReviews(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "reviews", ids)
}

// --- FAQ ---

func (r *ContentRepo) FAQ(ctx context.Context) ([]models.FAQ, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, question, answer FROM faq_items ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.FAQ
	for rows.Next() {
		var f models.FAQ
		if err := rows.Scan(&f.ID, &f.Question, &f.Answer); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreateFAQ(ctx context.Context, f models.FAQ, sortOrder int) error {
	ts := nowUTC()
	if f.ID == 0 {
		var next int
		_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(id), 0) + 1 FROM faq_items`).Scan(&next)
		f.ID = next
	}
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO faq_items (id, question, answer, sort_order, created_at, updated_at) VALUES (?,?,?,?,?,?)`),
		f.ID, f.Question, f.Answer, sortOrder, ts, ts)
	return err
}

func (r *ContentRepo) UpdateFAQ(ctx context.Context, f models.FAQ) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `UPDATE faq_items SET question=?, answer=?, updated_at=? WHERE id=?`),
		f.Question, f.Answer, ts, f.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeleteFAQ(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM faq_items WHERE id = ?`), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) ReorderFAQ(ctx context.Context, ids []int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, id := range ids {
		if _, err := tx.ExecContext(ctx, rebind(r.driver, `UPDATE faq_items SET sort_order = ? WHERE id = ?`), i, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *ContentRepo) NextFAQID(ctx context.Context) (int, error) {
	var next int
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(id), 0) + 1 FROM faq_items`).Scan(&next)
	return next, err
}

// --- Gallery ---

func (r *ContentRepo) Gallery(ctx context.Context) ([]models.GalleryItem, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, category, image_url, alt FROM gallery_items ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.GalleryItem
	for rows.Next() {
		var g models.GalleryItem
		if err := rows.Scan(&g.ID, &g.Title, &g.Category, &g.ImageURL, &g.Alt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreateGallery(ctx context.Context, g models.GalleryItem, sortOrder int) error {
	ts := nowUTC()
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO gallery_items (id, title, category, image_url, alt, sort_order, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`),
		g.ID, g.Title, g.Category, g.ImageURL, g.Alt, sortOrder, ts, ts)
	return err
}

func (r *ContentRepo) UpdateGallery(ctx context.Context, g models.GalleryItem) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `
UPDATE gallery_items SET title=?, category=?, image_url=?, alt=?, updated_at=? WHERE id=?`),
		g.Title, g.Category, g.ImageURL, g.Alt, ts, g.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) BackfillGalleryMedia(ctx context.Context) error {
	seedByID := map[string]models.GalleryItem{}
	for _, g := range defaultGallery() {
		seedByID[g.ID] = g
	}
	existing, err := r.Gallery(ctx)
	if err != nil {
		return err
	}
	for _, g := range existing {
		if g.ImageURL != "" && !containsStaleSeedMedia([]string{g.ImageURL}) {
			continue
		}
		seed, ok := seedByID[g.ID]
		if !ok {
			continue
		}
		g.ImageURL = seed.ImageURL
		if g.Alt == "" || containsStaleSeedMedia([]string{g.ImageURL}) {
			g.Alt = seed.Alt
		}
		if err := r.UpdateGallery(ctx, g); err != nil {
			return err
		}
	}
	return nil
}

func (r *ContentRepo) DeleteGallery(ctx context.Context, id string) error {
	return r.deleteByID(ctx, "gallery_items", id)
}

func (r *ContentRepo) ReorderGallery(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "gallery_items", ids)
}

// --- Nav ---

func (r *ContentRepo) Nav(ctx context.Context, visibleOnly bool) ([]models.NavLink, error) {
	q := `SELECT id, label, path, visible FROM nav_links`
	if visibleOnly {
		q += ` WHERE visible = 1`
	}
	q += ` ORDER BY sort_order`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.NavLink
	for rows.Next() {
		var n models.NavLink
		var vis int
		if err := rows.Scan(&n.ID, &n.Label, &n.Path, &vis); err != nil {
			return nil, err
		}
		n.Visible = vis == 1
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *ContentRepo) CreateNav(ctx context.Context, n models.NavLink, sortOrder int) error {
	ts := nowUTC()
	_, err := r.db.ExecContext(ctx, rebind(r.driver, `
INSERT INTO nav_links (id, label, path, sort_order, visible, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`),
		n.ID, n.Label, n.Path, sortOrder, boolInt(n.Visible), ts, ts)
	return err
}

func (r *ContentRepo) UpdateNav(ctx context.Context, n models.NavLink) error {
	ts := nowUTC()
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `
UPDATE nav_links SET label=?, path=?, visible=?, updated_at=? WHERE id=?`),
		n.Label, n.Path, boolInt(n.Visible), ts, n.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) DeleteNav(ctx context.Context, id string) error {
	return r.deleteByID(ctx, "nav_links", id)
}

func (r *ContentRepo) ReorderNav(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "nav_links", ids)
}

func (r *ContentRepo) deleteByID(ctx context.Context, table, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, table)), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (r *ContentRepo) MaxSort(ctx context.Context, table string) (int, error) {
	var max sql.NullInt64
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT MAX(sort_order) FROM %s`, table)).Scan(&max); err != nil {
		return 0, err
	}
	if !max.Valid {
		return 0, nil
	}
	return int(max.Int64) + 1, nil
}
