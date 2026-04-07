package migrations

import (
	"database/sql"

	"github.com/tss182/migration/internal/migrator"
	"github.com/tss182/migration/internal/schema"
)

func init() {
	migrator.RegisterMigration("20260407104034_add_users_table", &AddUsersTable{})
}

// AddUsersTable represents the 20260407104034_add_users_table migration.
type AddUsersTable struct{}

// Up applies the migration.
func (m *AddUsersTable) Up(db *sql.DB) error {
	return schema.CreateTable(db, migrator.Driver, "your_table", func(u *schema.TableBuilder) {
		u.Int("id").AutoIncrement().PrimaryKey()
		u.String("name").NotNull()
		u.Timestamps()
	})
}

// Down reverses the migration.
func (m *AddUsersTable) Down(db *sql.DB) error {
	return schema.DropTable(db, "your_table")
}
