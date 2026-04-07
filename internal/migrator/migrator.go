package migrator

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Driver is the active SQL driver ("mysql" or "postgres").
// Set this before calling Run, Rollback, or Status.
var Driver = "mysql"

// Migration is the interface every migration file must implement.
type Migration interface {
	Up(db *sql.DB) error   // Apply the migration
	Down(db *sql.DB) error // Reverse the migration
}

type migrationEntry struct {
	Name      string
	Migration Migration
}

var registry []migrationEntry

// RegisterMigration adds a migration to the global registry.
// Call this from an init() function inside each migration file.
func RegisterMigration(name string, m Migration) {
	registry = append(registry, migrationEntry{Name: name, Migration: m})
}

// Run executes all pending migrations in alphabetical (timestamp) order.
func Run(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to prepare migrations table: %w", err)
	}

	ran, err := getRanMigrations(db)
	if err != nil {
		return err
	}

	batch, err := getNextBatch(db)
	if err != nil {
		return err
	}

	sorted := sortedRegistry()
	count := 0

	for _, entry := range sorted {
		if ran[entry.Name] {
			continue
		}
		fmt.Printf("Migrating: %s\n", entry.Name)
		if err := entry.Migration.Up(db); err != nil {
			return fmt.Errorf("migration %q failed: %w", entry.Name, err)
		}
		if err := recordMigration(db, entry.Name, batch); err != nil {
			return fmt.Errorf("failed to record migration %q: %w", entry.Name, err)
		}
		fmt.Printf("Migrated:  %s\n", entry.Name)
		count++
	}

	if count == 0 {
		fmt.Println("Nothing to migrate.")
	} else {
		fmt.Printf("\n%d migration(s) ran successfully.\n", count)
	}
	return nil
}

// Rollback reverts the last `step` batches of migrations.
func Rollback(db *sql.DB, step int) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to prepare migrations table: %w", err)
	}

	toRollback, err := getMigrationsToRollback(db, step)
	if err != nil {
		return err
	}

	regMap := make(map[string]Migration, len(registry))
	for _, e := range registry {
		regMap[e.Name] = e.Migration
	}

	for _, name := range toRollback {
		m, ok := regMap[name]
		if !ok {
			return fmt.Errorf("migration %q not found in registry (did you forget to import it?)", name)
		}
		fmt.Printf("Rolling back: %s\n", name)
		if err := m.Down(db); err != nil {
			return fmt.Errorf("rollback of %q failed: %w", name, err)
		}
		if err := deleteMigration(db, name); err != nil {
			return fmt.Errorf("failed to remove record for %q: %w", name, err)
		}
		fmt.Printf("Rolled back:  %s\n", name)
	}

	if len(toRollback) == 0 {
		fmt.Println("Nothing to rollback.")
	} else {
		fmt.Printf("\n%d migration(s) rolled back.\n", len(toRollback))
	}
	return nil
}

// Status prints a table showing every registered migration and whether it has run.
func Status(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to prepare migrations table: %w", err)
	}

	ran, err := getRanMigrationsWithBatch(db)
	if err != nil {
		return err
	}

	fmt.Printf("%-60s %-8s %s\n", "Migration", "Batch", "Status")
	fmt.Println(strings.Repeat("-", 80))

	for _, entry := range sortedRegistry() {
		if batch, ok := ran[entry.Name]; ok {
			fmt.Printf("%-60s %-8d %s\n", entry.Name, batch, "Ran")
		} else {
			fmt.Printf("%-60s %-8s %s\n", entry.Name, "-", "Pending")
		}
	}
	return nil
}

// ----- internal helpers -----

func sortedRegistry() []migrationEntry {
	cp := make([]migrationEntry, len(registry))
	copy(cp, registry)
	sort.Slice(cp, func(i, j int) bool { return cp[i].Name < cp[j].Name })
	return cp
}

// ph returns the placeholder token for the current driver (? or $N).
func ph(n int) string {
	if Driver == "postgres" {
		return "$" + strconv.Itoa(n)
	}
	return "?"
}

func ensureMigrationsTable(db *sql.DB) error {
	var q string
	if Driver == "postgres" {
		q = `CREATE TABLE IF NOT EXISTS migrations (
			id       SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL,
			batch    INTEGER NOT NULL
		)`
	} else {
		q = `CREATE TABLE IF NOT EXISTS migrations (
			id        INT UNSIGNED NOT NULL AUTO_INCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch     INT NOT NULL,
			PRIMARY KEY (id)
		)`
	}
	_, err := db.Exec(q)
	return err
}

func getRanMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT migration FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ran := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		ran[name] = true
	}
	return ran, rows.Err()
}

func getRanMigrationsWithBatch(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query("SELECT migration, batch FROM migrations ORDER BY batch, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ran := make(map[string]int)
	for rows.Next() {
		var name string
		var batch int
		if err := rows.Scan(&name, &batch); err != nil {
			return nil, err
		}
		ran[name] = batch
	}
	return ran, rows.Err()
}

func getNextBatch(db *sql.DB) (int, error) {
	var maxBatch sql.NullInt64
	if err := db.QueryRow("SELECT MAX(batch) FROM migrations").Scan(&maxBatch); err != nil || !maxBatch.Valid {
		return 1, nil
	}
	return int(maxBatch.Int64) + 1, nil
}

func recordMigration(db *sql.DB, name string, batch int) error {
	q := fmt.Sprintf("INSERT INTO migrations (migration, batch) VALUES (%s, %s)", ph(1), ph(2))
	_, err := db.Exec(q, name, batch)
	return err
}

func deleteMigration(db *sql.DB, name string) error {
	q := fmt.Sprintf("DELETE FROM migrations WHERE migration = %s", ph(1))
	_, err := db.Exec(q, name)
	return err
}

func getMigrationsToRollback(db *sql.DB, steps int) ([]string, error) {
	var maxBatch sql.NullInt64
	if err := db.QueryRow("SELECT MAX(batch) FROM migrations").Scan(&maxBatch); err != nil || !maxBatch.Valid {
		return nil, nil
	}

	minBatch := int(maxBatch.Int64) - steps + 1
	if minBatch < 1 {
		minBatch = 1
	}

	q := fmt.Sprintf(
		"SELECT migration FROM migrations WHERE batch >= %s ORDER BY batch DESC, id DESC",
		ph(1),
	)
	rows, err := db.Query(q, minBatch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}
