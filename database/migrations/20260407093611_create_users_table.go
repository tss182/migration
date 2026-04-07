package migrations

import (
	"database/sql"

	"github.com/tss182/migration/internal/migrator"
)

func init() {
	migrator.RegisterMigration("20260407093611_create_users_table", &CreateUsersTable{})
}

// CreateUsersTable represents the 20260407093611_create_users_table migration.
type CreateUsersTable struct{}

// Up applies the migration.
func (m *CreateUsersTable) Up(db *sql.DB) error {
	_, err := db.Exec(`
		-- TODO: write your migration SQL here
	`)
	return err
}

// Down reverses the migration.
func (m *CreateUsersTable) Down(db *sql.DB) error {
	_, err := db.Exec(`
		-- TODO: write your rollback SQL here
	`)
	return err
}
