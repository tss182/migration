package config

import (
	"fmt"
	"os"
)

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// Load reads database configuration from environment variables.
func Load() DBConfig {
	return DBConfig{
		Driver:   getEnv("DB_DRIVER", "mysql"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		Database: getEnv("DB_DATABASE", ""),
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", ""),
	}
}

// DSN builds the data source name for the configured driver.
func (c DBConfig) DSN() (string, error) {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
			c.Username, c.Password, c.Host, c.Port, c.Database), nil
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.Username, c.Password, c.Database), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", c.Driver)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
