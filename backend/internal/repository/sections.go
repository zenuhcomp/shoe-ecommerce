package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

const sectionSelectCols = `
id, type, enabled, sort_order, title, subtitle, badge, content,
background_type, background_color, background_image, background_video, overlay_opacity,
animation, animation_duration, animation_delay, hover_effect,
primary_cta_text, primary_cta_url, secondary_cta_text, secondary_cta_url, extra_data`

func scanSection(scanner interface {
	Scan(dest ...any) error
}) (models.Section, error) {
	var s models.Section
	var enabled int
	if err := scanner.Scan(
		&s.ID, &s.Type, &enabled, &s.SortOrder, &s.Title, &s.Subtitle, &s.Badge, &s.Content,
		&s.BackgroundType, &s.BackgroundColor, &s.BackgroundImage, &s.BackgroundVideo, &s.OverlayOpacity,
		&s.Animation, &s.AnimationDuration, &s.AnimationDelay, &s.HoverEffect,
		&s.PrimaryCtaText, &s.PrimaryCtaUrl, &s.SecondaryCtaText, &s.SecondaryCtaUrl, &s.ExtraData,
	); err != nil {
		return s, err
	}
	s.Enabled = enabled == 1
	if s.Animation == "" {
		s.Animation = "fade-up"
	}
	if s.AnimationDuration == "" {
		s.AnimationDuration = "0.8s"
	}
	if s.AnimationDelay == "" {
		s.AnimationDelay = "0.1s"
	}
	if s.HoverEffect == "" {
		s.HoverEffect = "lift"
	}
	if s.BackgroundType == "" {
		s.BackgroundType = "default"
	}
	return s, nil
}

func (r *ContentRepo) Sections(ctx context.Context, includeDisabled bool) ([]models.Section, error) {
	q := fmt.Sprintf(`SELECT %s FROM sections`, sectionSelectCols)
	if !includeDisabled {
		q += ` WHERE enabled = 1`
	}
	q += ` ORDER BY sort_order ASC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Section
	for rows.Next() {
		s, err := scanSection(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ContentRepo) SectionByID(ctx context.Context, id string) (models.Section, error) {
	q := rebind(r.driver, fmt.Sprintf(`SELECT %s FROM sections WHERE id = ?`, sectionSelectCols))
	s, err := scanSection(r.db.QueryRowContext(ctx, q, id))
	if err == sql.ErrNoRows {
		return s, fmt.Errorf("section not found")
	}
	return s, err
}

func (r *ContentRepo) CreateSection(ctx context.Context, s models.Section, sortOrder int) error {
	ts := nowUTC()
	q := rebind(r.driver, `
INSERT INTO sections (
	id, type, enabled, sort_order, title, subtitle, badge, content,
	background_type, background_color, background_image, background_video, overlay_opacity,
	animation, animation_duration, animation_delay, hover_effect,
	primary_cta_text, primary_cta_url, secondary_cta_text, secondary_cta_url, extra_data,
	created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := r.db.ExecContext(ctx, q,
		s.ID, s.Type, boolInt(s.Enabled), sortOrder, s.Title, s.Subtitle, s.Badge, s.Content,
		s.BackgroundType, s.BackgroundColor, s.BackgroundImage, s.BackgroundVideo, s.OverlayOpacity,
		s.Animation, s.AnimationDuration, s.AnimationDelay, s.HoverEffect,
		s.PrimaryCtaText, s.PrimaryCtaUrl, s.SecondaryCtaText, s.SecondaryCtaUrl, s.ExtraData,
		ts, ts,
	)
	return err
}

func (r *ContentRepo) UpdateSection(ctx context.Context, s models.Section) error {
	ts := nowUTC()
	q := rebind(r.driver, `
UPDATE sections SET
	type=?, enabled=?, title=?, subtitle=?, badge=?, content=?,
	background_type=?, background_color=?, background_image=?, background_video=?, overlay_opacity=?,
	animation=?, animation_duration=?, animation_delay=?, hover_effect=?,
	primary_cta_text=?, primary_cta_url=?, secondary_cta_text=?, secondary_cta_url=?, extra_data=?,
	updated_at=?
WHERE id=?`)
	res, err := r.db.ExecContext(ctx, q,
		s.Type, boolInt(s.Enabled), s.Title, s.Subtitle, s.Badge, s.Content,
		s.BackgroundType, s.BackgroundColor, s.BackgroundImage, s.BackgroundVideo, s.OverlayOpacity,
		s.Animation, s.AnimationDuration, s.AnimationDelay, s.HoverEffect,
		s.PrimaryCtaText, s.PrimaryCtaUrl, s.SecondaryCtaText, s.SecondaryCtaUrl, s.ExtraData,
		ts, s.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("section not found")
	}
	return nil
}

func (r *ContentRepo) DeleteSection(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, rebind(r.driver, `DELETE FROM sections WHERE id = ?`), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("section not found")
	}
	return nil
}

func (r *ContentRepo) DeleteAllSections(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sections`)
	return err
}

func (r *ContentRepo) ReorderSections(ctx context.Context, ids []string) error {
	return r.reorderStrings(ctx, "sections", ids)
}
