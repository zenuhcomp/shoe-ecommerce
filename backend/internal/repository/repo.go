package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kibsoft/luna-collections/backend/internal/models"
)

type ContentRepo struct {
	db     *sql.DB
	driver string
}

func NewContentRepo(db *sql.DB, driver string) *ContentRepo {
	return &ContentRepo{db: db, driver: driver}
}

func (r *ContentRepo) Driver() string { return r.driver }

func (r *ContentRepo) SeedIfEmpty(ctx context.Context) error {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM site_config`).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		if err := r.seedAll(ctx); err != nil {
			return err
		}
	}

	// Check sections
	var sectionCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sections`).Scan(&sectionCount); err == nil && sectionCount == 0 {
		_ = r.applyFootwearPreset(ctx, false)
	}

	// Check media
	_ = r.SeedDefaultMediaIfEmpty(ctx)

	return nil
}

func (r *ContentRepo) seedAll(ctx context.Context) error {
	if err := r.UpdateSite(ctx, defaultSiteConfig()); err != nil {
		return err
	}
	for i, s := range defaultServices() {
		svc := s
		svc.Category = s.Style
		svc.Badge = "Official Wear"
		svc.Currency = "KES"
		svc.InStock = true
		svc.Specs = map[string]string{
			"Material":    "Full-Grain Leather",
			"Heel Height": s.Length,
			"Occasion":    "Boardroom & Office",
		}
		if err := r.CreateService(ctx, svc, i); err != nil {
			return err
		}
	}
	for i, f := range defaultFeatures() {
		if err := r.CreateFeature(ctx, f, i); err != nil {
			return err
		}
	}
	for i, p := range defaultPress() {
		if err := r.CreatePress(ctx, p, i); err != nil {
			return err
		}
	}
	for i, rev := range defaultReviews() {
		if err := r.CreateReview(ctx, rev, i); err != nil {
			return err
		}
	}
	for i, item := range defaultFAQ() {
		if err := r.CreateFAQ(ctx, item, i); err != nil {
			return err
		}
	}
	for i, g := range defaultGallery() {
		if err := r.CreateGallery(ctx, g, i); err != nil {
			return err
		}
	}
	for i, link := range defaultNav() {
		if err := r.CreateNav(ctx, link, i); err != nil {
			return err
		}
	}
	return nil
}

func (r *ContentRepo) Home(ctx context.Context) (models.HomePage, error) {
	site, err := r.Site(ctx)
	if err != nil {
		return models.HomePage{}, err
	}
	services, err := r.Services(ctx)
	if err != nil {
		return models.HomePage{}, err
	}
	features, err := r.Features(ctx)
	if err != nil {
		return models.HomePage{}, err
	}
	press, err := r.Press(ctx)
	if err != nil {
		return models.HomePage{}, err
	}
	reviews, err := r.Reviews(ctx, true)
	if err != nil {
		return models.HomePage{}, err
	}
	faq, err := r.FAQ(ctx)
	if err != nil {
		return models.HomePage{}, err
	}
	nav, err := r.Nav(ctx, true)
	if err != nil {
		return models.HomePage{}, err
	}
	sections, err := r.Sections(ctx, false)
	if err != nil {
		sections = []models.Section{}
	}

	return models.HomePage{
		Site:     site,
		Services: services,
		Features: features,
		Press:    press,
		Reviews:  reviews,
		FAQ:      faq,
		Nav:      nav,
		Sections: sections,
	}, nil
}

func nowUTC() time.Time { return time.Now().UTC() }

func (r *ContentRepo) isEmpty(ctx context.Context, table string) (bool, error) {
	var n int
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)
	if err := r.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return false, err
	}
	return n == 0, nil
}
