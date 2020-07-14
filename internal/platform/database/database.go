package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // to register postgres driver
)

// Config is the configuration that is required for the database connection
type Config struct {
	User string
	Password string
	Host string
	Path string
	DisableSSL bool
}

// Open will open a connection to database based on configuration
func Open(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "require")
	if cfg.DisableSSL {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Path,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}