package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// ModeTest uses SQLite; ModeProduction uses PostgreSQL.
const (
	ModeTest       = "test"
	ModeProduction = "production"
)

type Config struct {
	AppEnv         string
	BackendPort    string
	FrontendPort   string
	DatabaseMode   string
	DatabaseURL    string
	SQLitePath     string
	Postgres       PostgresConfig
	CORSOrigins    []string
	JWTSecret      string
	AdminEmail     string
	AdminPassword  string
	AdminName      string
	JWTExpiryHours int

	// Safaricom M-Pesa Configuration (Production / .env)
	SafaricomBaseURL string
	ConsumerKey      string
	ConsumerSecret   string
	PassKey          string
	ShortCode        string
	CallbackURL      string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SSLMode  string
}

func Load() (Config, error) {
	envFile := defaultEnvPath()
	_ = godotenv.Overload(envFile)

	cfg := Config{
		AppEnv:       getenv("APP_ENV", "development"),
		BackendPort:  getenv("BACKEND_PORT", getenv("PORT", "8083")),
		FrontendPort: getenv("FRONTEND_PORT", "4203"),
		DatabaseMode: strings.ToLower(getenv("DATABASE_MODE", ModeTest)),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		SQLitePath:   getenv("SQLITE_PATH", "./data/luna_collections.db"),
		Postgres: PostgresConfig{
			Host:     getenv("POSTGRES_HOST", "localhost"),
			Port:     getenv("POSTGRES_PORT", "5432"),
			User:     getenv("POSTGRES_USER", "luna"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       getenv("POSTGRES_DB", "luna_collections"),
			SSLMode:  getenv("POSTGRES_SSLMODE", "disable"),
		},
		JWTSecret:      getenv("JWT_SECRET", "change-me-in-production-use-long-random-string"),
		AdminEmail:     getenv("CMS_ADMIN_EMAIL", "admin@lunacollections.local"),
		AdminPassword:  getenv("CMS_ADMIN_PASSWORD", "luna-admin-change-me"),
		AdminName:      getenv("CMS_ADMIN_NAME", "Site Admin"),
		JWTExpiryHours: getenvInt("JWT_EXPIRY_HOURS", 24),

		SafaricomBaseURL: getenv("SAFARICOM_BASE_URL", "https://api.safaricom.co.ke"),
		ConsumerKey:      os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:   os.Getenv("CONSUMER_SECRET"),
		PassKey:          os.Getenv("PASSKEY"),
		ShortCode:        firstNonEmpty(os.Getenv("COUNTY_SHORTC0DE"), os.Getenv("COUNTY_SHORTCODE"), os.Getenv("SHORT_CODE"), os.Getenv("SHORTCODE"), os.Getenv("BUSINESS_SHORTCODE")),
		CallbackURL:      os.Getenv("MPESA_CALLBACK_URL"),
	}

	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		cfg.CORSOrigins = splitCSV(raw)
	} else {
		cfg.CORSOrigins = []string{
			fmt.Sprintf("http://localhost:%s", cfg.FrontendPort),
			fmt.Sprintf("http://127.0.0.1:%s", cfg.FrontendPort),
		}
	}

	return cfg, cfg.Validate()
}

func defaultEnvPath() string {
	if p := os.Getenv("ENV_FILE"); p != "" {
		return p
	}
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			c1 := filepath.Join(dir, ".env")
			if fi, err := os.Stat(c1); err == nil && !fi.IsDir() {
				return c1
			}
			c2 := filepath.Join(dir, "backend", ".env")
			if fi, err := os.Stat(c2); err == nil && !fi.IsDir() {
				return c2
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ".env"
}

func (c Config) Validate() error {
	switch c.DatabaseMode {
	case ModeTest, ModeProduction:
	default:
		return fmt.Errorf("DATABASE_MODE must be %q or %q, got %q", ModeTest, ModeProduction, c.DatabaseMode)
	}
	if c.DatabaseMode == ModeProduction && c.DatabaseURL == "" {
		if c.Postgres.DB == "" || c.Postgres.User == "" {
			return fmt.Errorf("production mode requires POSTGRES_DB and POSTGRES_USER (or DATABASE_URL)")
		}
	}
	return nil
}

func (c Config) Driver() string {
	if c.DatabaseMode == ModeProduction {
		return "postgres"
	}
	return "sqlite"
}

func (c Config) PostgresDSN() string {
	if c.DatabaseURL != "" && strings.HasPrefix(c.DatabaseURL, "postgres") {
		return c.DatabaseURL
	}
	p := c.Postgres
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DB, p.SSLMode,
	)
}

func (c Config) SQLiteDSN() string {
	if c.DatabaseURL != "" && (strings.HasPrefix(c.DatabaseURL, "file:") || strings.HasSuffix(c.DatabaseURL, ".db")) {
		return c.DatabaseURL
	}
	path := c.SQLitePath
	if !strings.Contains(path, "?") {
		path += "?_pragma=foreign_keys(1)"
	}
	return path
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func firstNonEmpty(items ...string) string {
	for _, s := range items {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}
