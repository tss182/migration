package seeders

import (
	"database/sql"

	"github.com/tss182/migration/internal/migrator"
)

func init() {
	migrator.RegisterSeeder("UserSeeder", &UserSeeder{})
}

// UserSeeder is a database seeder.
type UserSeeder struct{}

// Run executes the seeder.
func (s *UserSeeder) Run(db *sql.DB) error {
	_, err := db.Exec(`
		-- TODO: write your seed SQL here
	`)
	return err
}

// Down reverses the seeder.
func (s *UserSeeder) Down(db *sql.DB) error {
	_, err := db.Exec(`
		-- TODO: write your seed rollback SQL here
	`)
	return err
}
