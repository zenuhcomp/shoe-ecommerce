package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSiteSeedsMissingSingletonRow(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE site_config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		data TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create site_config: %v", err)
	}

	repo := NewContentRepo(db, "sqlite")
	site, err := repo.Site(context.Background())
	if err != nil {
		t.Fatalf("Site() error = %v", err)
	}
	if site.BrandName == "" {
		t.Fatal("Site() returned an empty default site")
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM site_config WHERE id = 1`).Scan(&count); err != nil {
		t.Fatalf("count seeded site config: %v", err)
	}
	if count != 1 {
		t.Fatalf("seeded site_config rows = %d, want 1", count)
	}
}
