package migrator

import (
	"database/sql"
	"fmt"
)

// Seeder is the interface every seeder file must implement.
type Seeder interface {
	Run(db *sql.DB) error
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

// RunSeeders executes all registered seeders in registration order.
func RunSeeders(db *sql.DB) error {
	if len(seederRegistry) == 0 {
		fmt.Println("No seeders registered.")
		return nil
	}

	for _, entry := range seederRegistry {
		fmt.Printf("Seeding: %s\n", entry.Name)
		if err := entry.Seeder.Run(db); err != nil {
			return fmt.Errorf("seeder %q failed: %w", entry.Name, err)
		}
		fmt.Printf("Seeded:  %s\n", entry.Name)
	}

	fmt.Printf("\n%d seeder(s) ran successfully.\n", len(seederRegistry))
	return nil
}
