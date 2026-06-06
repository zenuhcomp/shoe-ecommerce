package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/kibsoft/luna-collections/backend/internal/config"
)

type DB struct {
	*sql.DB
	Driver string
	Mode   string
}

func Connect(ctx context.Context, cfg config.Config) (*DB, error) {
	driver := cfg.Driver()
	var dsn string
	var openDriver string
	switch driver {
	case "postgres":
		dsn = cfg.PostgresDSN()
		openDriver = "pgx"
	case "sqlite":
		dsn = cfg.SQLiteDSN()
		openDriver = "sqlite"
	default:
		return nil, fmt.Errorf("unsupported driver %q", driver)
	}

	sqlDB, err := sql.Open(openDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	wrapped := &DB{DB: sqlDB, Driver: openDriver, Mode: cfg.DatabaseMode}
	if err := wrapped.migrate(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	log.Printf("database connected (%s, mode=%s)", driver, cfg.DatabaseMode)
	return wrapped, nil
}

