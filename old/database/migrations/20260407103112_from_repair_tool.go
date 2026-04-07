package migrations

import (
	"database/sql"

	"github.com/tss182/migration/internal/migrator"
	"github.com/tss182/migration/internal/schema"
)

func init() {
	migrator.RegisterMigration("20260407103112_from_repair_tool", &FromRepairTool{})
}

// FromRepairTool represents the 20260407103112_from_repair_tool migration.
type FromRepairTool struct{}

// Up applies the migration.
func (m *FromRepairTool) Up(db *sql.DB) error {
	return schema.CreateTable(db, migrator.Driver, "your_table", func(u *schema.TableBuilder) {
		u.Int("id").AutoIncrement().PrimaryKey()
		u.String("name").NotNull()
		u.Timestamps()
	})
}

// Down reverses the migration.
func (m *FromRepairTool) Down(db *sql.DB) error {
	return schema.DropTable(db, "your_table")
}
