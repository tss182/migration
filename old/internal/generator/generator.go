package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/tss182/migration/internal/project"
	"github.com/tss182/migration/internal/stubs"
)

// CreateMigration ensures project folders exist and writes a new migration file.
func CreateMigration(baseDir, name string) (string, error) {
	if err := project.EnsureRegistryPackages(baseDir); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102150405")
	migrationName := timestamp + "_" + name
	structName := toPascalCase(name)
	filename := migrationName + ".go"

	stub, err := stubs.FS.ReadFile("migration.stub")
	if err != nil {
		return "", fmt.Errorf("failed to read migration stub: %w", err)
	}

	content := strings.ReplaceAll(string(stub), "{{NAME}}", migrationName)
	content = strings.ReplaceAll(content, "{{STRUCT}}", structName)

	dir := filepath.Join(baseDir, "database", "migrations")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return path, nil
}

// CreateSeeder ensures project folders exist and writes a new seeder file.
func CreateSeeder(baseDir, name string) (string, error) {
	if err := project.EnsureRegistryPackages(baseDir); err != nil {
		return "", err
	}

	structName := toPascalCase(name)
	filename := time.Now().Format("20060102150405") + "_" + strings.ToLower(name) + ".go"

	stub, err := stubs.FS.ReadFile("seeder.stub")
	if err != nil {
		return "", fmt.Errorf("failed to read seeder stub: %w", err)
	}

	content := strings.ReplaceAll(string(stub), "{{NAME}}", structName)
	content = strings.ReplaceAll(content, "{{STRUCT}}", structName)

	dir := filepath.Join(baseDir, "database", "seeders")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return path, nil
}
