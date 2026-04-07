package db

import (
	"database/sql"
	"fmt"

	"github.com/tss182/migration/internal/config"

	// Register MySQL driver
	_ "github.com/go-sql-driver/mysql"
	// Register PostgreSQL driver
	_ "github.com/lib/pq"
)

// Connect opens and verifies a database connection using the given config.
func Connect(cfg config.DBConfig) (*sql.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to %s at %s:%s — %w", cfg.Driver, cfg.Host, cfg.Port, err)
	}

	return conn, nil
}
