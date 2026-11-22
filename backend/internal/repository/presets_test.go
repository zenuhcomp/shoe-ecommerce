package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, *ContentRepo) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stmts := []string{
		`CREATE TABLE site_config (id INTEGER PRIMARY KEY CHECK (id = 1), data TEXT NOT NULL, updated_at TIMESTAMP NOT NULL)`,
		`CREATE TABLE lash_services (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, slug TEXT NOT NULL UNIQUE, emoji TEXT NOT NULL DEFAULT '',
			style TEXT NOT NULL DEFAULT '', length TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '',
			price_kes INTEGER NOT NULL DEFAULT 0, featured INTEGER NOT NULL DEFAULT 0, images TEXT NOT NULL DEFAULT '[]',
			view_360 TEXT NOT NULL DEFAULT '[]', category TEXT NOT NULL DEFAULT '', badge TEXT NOT NULL DEFAULT '',
			original_price_kes INTEGER NOT NULL DEFAULT 0, currency TEXT NOT NULL DEFAULT 'KES', specs TEXT NOT NULL DEFAULT '{}',
			video_url TEXT NOT NULL DEFAULT '', in_stock INTEGER NOT NULL DEFAULT 1, sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE sections (
			id TEXT PRIMARY KEY, type TEXT NOT NULL, enabled INTEGER NOT NULL DEFAULT 1, sort_order INTEGER NOT NULL DEFAULT 0,
			title TEXT NOT NULL DEFAULT '', subtitle TEXT NOT NULL DEFAULT '', badge TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '', background_type TEXT NOT NULL DEFAULT 'default', background_color TEXT NOT NULL DEFAULT '',
			background_image TEXT NOT NULL DEFAULT '', background_video TEXT NOT NULL DEFAULT '', overlay_opacity INTEGER NOT NULL DEFAULT 40,
			animation TEXT NOT NULL DEFAULT 'fade-up', animation_duration TEXT NOT NULL DEFAULT '0.8s', animation_delay TEXT NOT NULL DEFAULT '0.1s',
			hover_effect TEXT NOT NULL DEFAULT 'lift', primary_cta_text TEXT NOT NULL DEFAULT '', primary_cta_url TEXT NOT NULL DEFAULT '',
			secondary_cta_text TEXT NOT NULL DEFAULT '', secondary_cta_url TEXT NOT NULL DEFAULT '', extra_data TEXT NOT NULL DEFAULT '{}',
			created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE media_assets (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, category TEXT NOT NULL DEFAULT 'general', media_type TEXT NOT NULL DEFAULT 'image',
			url TEXT NOT NULL, thumbnail_url TEXT NOT NULL DEFAULT '', alt_text TEXT NOT NULL DEFAULT '', tags TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("setup stmt failed: %v", err)
		}
	}
	repo := NewContentRepo(db, "sqlite")
	return db, repo
}

func TestPresetsApplication(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// 1. List presets
	presets := repo.ListPresets()
	if len(presets) < 4 {
		t.Fatalf("expected at least 4 presets, got %d", len(presets))
	}

	// 2. Apply furniture preset
	if err := repo.ApplyPreset(ctx, "furniture", true); err != nil {
		t.Fatalf("apply furniture preset: %v", err)
	}
	site, err := repo.Site(ctx)
	if err != nil {
		t.Fatalf("get site: %v", err)
	}
	if site.Industry != "furniture" {
		t.Fatalf("expected industry furniture, got %s", site.Industry)
	}
	services, err := repo.Services(ctx)
	if err != nil {
		t.Fatalf("get services: %v", err)
	}
	if len(services) == 0 {
		t.Fatal("expected seeded furniture products")
	}
	if len(services[0].Specs) == 0 {
		t.Fatal("expected specs in furniture product")
	}

	// 3. Apply vehicles preset
	if err := repo.ApplyPreset(ctx, "vehicles", true); err != nil {
		t.Fatalf("apply vehicles preset: %v", err)
	}
	site, err = repo.Site(ctx)
	if err != nil {
		t.Fatalf("get site: %v", err)
	}
	if site.Industry != "vehicles" {
		t.Fatalf("expected industry vehicles, got %s", site.Industry)
	}

	// Check sections
	sections, err := repo.Sections(ctx, false)
	if err != nil {
		t.Fatalf("get sections: %v", err)
	}
	if len(sections) == 0 {
		t.Fatal("expected active sections")
	}
}
