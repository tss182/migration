package project

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	migrationsRegister = "// Package migrations contains all database migration files.\npackage migrations\n"
	seedersRegister    = "// Package seeders contains all database seeder files.\npackage seeders\n"
)

// EnsureRegistryPackages makes sure migration/seeder package directories and register files exist.
func EnsureRegistryPackages(baseDir string) error {
	if err := ensureDirAndFile(baseDir, filepath.Join("database", "migrations"), "register.go", migrationsRegister); err != nil {
		return err
	}
	if err := ensureDirAndFile(baseDir, filepath.Join("database", "seeders"), "register.go", seedersRegister); err != nil {
		return err
	}
	return nil
}

func ensureDirAndFile(baseDir, dir, file, content string) error {
	absDir := filepath.Join(baseDir, dir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", absDir, err)
	}
	absFile := filepath.Join(absDir, file)
	if _, err := os.Stat(absFile); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s: %w", absFile, err)
	}
	if err := os.WriteFile(absFile, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to create %s: %w", absFile, err)
	}
	return nil
}
