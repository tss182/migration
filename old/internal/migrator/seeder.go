package migrator

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

// Seeder is the interface every seeder file must implement.
type Seeder interface {
	Run(db *sql.DB) error
}

// ReversibleSeeder is implemented by seeders that support rollback.
type ReversibleSeeder interface {
	Down(db *sql.DB) error
}

type seederEntry struct {
	Name   string
	Seeder Seeder
}

var seederRegistry []seederEntry

// RegisterSeeder adds a seeder to the global registry.
// Call this from an init() function inside each seeder file.
func RegisterSeeder(name string, s Seeder) {
	seederRegistry = append(seederRegistry, seederEntry{Name: name, Seeder: s})
}

// RunSeeders executes all pending seeders in registration order.
func RunSeeders(db *sql.DB) error {
	if err := ensureSeedersTable(db); err != nil {
		return fmt.Errorf("failed to prepare seeders table: %w", err)
	}

	if len(seederRegistry) == 0 {
		fmt.Println("No seeders registered.")
		return nil
	}

	ran, err := getRanSeeders(db)
	if err != nil {
		return err
	}

	batch, err := getNextSeederBatch(db)
	if err != nil {
		return err
	}

	sorted := sortedSeeders()
	count := 0

	for _, entry := range sorted {
		if ran[entry.Name] {
			continue
		}
		fmt.Printf("Seeding: %s\n", entry.Name)
		if err := entry.Seeder.Run(db); err != nil {
			return fmt.Errorf("seeder %q failed: %w", entry.Name, err)
		}
		if err := recordSeeder(db, entry.Name, batch); err != nil {
			return fmt.Errorf("failed to record seeder %q: %w", entry.Name, err)
		}
		fmt.Printf("Seeded:  %s\n", entry.Name)
		count++
	}

	if count == 0 {
		fmt.Println("Nothing to seed.")
	} else {
		fmt.Printf("\n%d seeder(s) ran successfully.\n", count)
	}
	return nil
}

// RollbackSeeders reverts the last `step` batches of seeders.
func RollbackSeeders(db *sql.DB, step int) error {
	if err := ensureSeedersTable(db); err != nil {
		return fmt.Errorf("failed to prepare seeders table: %w", err)
	}

	toRollback, err := getSeedersToRollback(db, step)
	if err != nil {
		return err
	}

	regMap := make(map[string]Seeder, len(seederRegistry))
	for _, e := range seederRegistry {
		regMap[e.Name] = e.Seeder
	}

	for _, name := range toRollback {
		s, ok := regMap[name]
		if !ok {
			return fmt.Errorf("seeder %q not found in registry (did you forget to import it?)", name)
		}
		reversible, ok := s.(ReversibleSeeder)
		if !ok {
			return fmt.Errorf("seeder %q does not implement Down(db) error", name)
		}
		fmt.Printf("Rolling back seeder: %s\n", name)
		if err := reversible.Down(db); err != nil {
			return fmt.Errorf("rollback of seeder %q failed: %w", name, err)
		}
		if err := deleteSeeder(db, name); err != nil {
			return fmt.Errorf("failed to remove record for seeder %q: %w", name, err)
		}
		fmt.Printf("Rolled back seeder:  %s\n", name)
	}

	if len(toRollback) == 0 {
		fmt.Println("Nothing to rollback.")
	} else {
		fmt.Printf("\n%d seeder(s) rolled back.\n", len(toRollback))
	}
	return nil
}

func sortedSeeders() []seederEntry {
	cp := make([]seederEntry, len(seederRegistry))
	copy(cp, seederRegistry)
	sort.Slice(cp, func(i, j int) bool { return cp[i].Name < cp[j].Name })
	return cp
}

func ensureSeedersTable(db *sql.DB) error {
	var q string
	if Driver == "postgres" {
		q = `CREATE TABLE IF NOT EXISTS seeders (
			id       SERIAL PRIMARY KEY,
			seeder   VARCHAR(255) NOT NULL,
			batch    INTEGER NOT NULL
		)`
	} else {
		q = `CREATE TABLE IF NOT EXISTS seeders (
			id       INT UNSIGNED NOT NULL AUTO_INCREMENT,
			seeder   VARCHAR(255) NOT NULL,
			batch    INT NOT NULL,
			PRIMARY KEY (id)
		)`
	}
	_, err := db.Exec(q)
	return err
}

func getRanSeeders(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT seeder FROM seeders")
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

func getNextSeederBatch(db *sql.DB) (int, error) {
	var maxBatch sql.NullInt64
	if err := db.QueryRow("SELECT MAX(batch) FROM seeders").Scan(&maxBatch); err != nil || !maxBatch.Valid {
		return 1, nil
	}
	return int(maxBatch.Int64) + 1, nil
}

func recordSeeder(db *sql.DB, name string, batch int) error {
	q := fmt.Sprintf("INSERT INTO seeders (seeder, batch) VALUES (%s, %s)", ph(1), ph(2))
	_, err := db.Exec(q, name, batch)
	return err
}

func deleteSeeder(db *sql.DB, name string) error {
	q := fmt.Sprintf("DELETE FROM seeders WHERE seeder = %s", ph(1))
	_, err := db.Exec(q, name)
	return err
}

func getSeedersToRollback(db *sql.DB, steps int) ([]string, error) {
	var maxBatch sql.NullInt64
	if err := db.QueryRow("SELECT MAX(batch) FROM seeders").Scan(&maxBatch); err != nil || !maxBatch.Valid {
		return nil, nil
	}

	minBatch := int(maxBatch.Int64) - steps + 1
	if minBatch < 1 {
		minBatch = 1
	}

	q := fmt.Sprintf(
		"SELECT seeder FROM seeders WHERE batch >= %s ORDER BY batch DESC, id DESC",
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

// SeederStatus prints a table showing every registered seeder and whether it has run.
func SeederStatus(db *sql.DB) error {
	if err := ensureSeedersTable(db); err != nil {
		return fmt.Errorf("failed to prepare seeders table: %w", err)
	}

	ran, err := getRanSeedersWithBatch(db)
	if err != nil {
		return err
	}

	fmt.Printf("%-60s %-8s %s\n", "Seeder", "Batch", "Status")
	fmt.Println(strings.Repeat("-", 80))

	for _, entry := range sortedSeeders() {
		if batch, ok := ran[entry.Name]; ok {
			fmt.Printf("%-60s %-8d %s\n", entry.Name, batch, "Ran")
		} else {
			fmt.Printf("%-60s %-8s %s\n", entry.Name, "-", "Pending")
		}
	}
	return nil
}

func getRanSeedersWithBatch(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query("SELECT seeder, batch FROM seeders ORDER BY batch, id")
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
