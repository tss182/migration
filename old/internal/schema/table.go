package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

// TableBuilder provides a fluent API for building CREATE TABLE SQL.
type TableBuilder struct {
	driver  string
	table   string
	columns []*ColumnBuilder
}

// ColumnBuilder provides chainable modifiers for a single column.
type ColumnBuilder struct {
	name          string
	typeName      string
	nullable      bool
	defaultSet    bool
	defaultValue  string
	defaultIsRaw  bool
	autoIncrement bool
	primaryKey    bool
	unique        bool
}

// CreateTable builds and executes a CREATE TABLE statement.
func CreateTable(db *sql.DB, driver, table string, build func(t *TableBuilder)) error {
	tb := NewTable(driver, table)
	build(tb)
	_, err := db.Exec(tb.BuildCreateSQL())
	return err
}

// DropTable executes DROP TABLE IF EXISTS.
func DropTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	return err
}

// NewTable initializes a fluent table builder.
func NewTable(driver, table string) *TableBuilder {
	return &TableBuilder{driver: driver, table: table}
}

// Int adds an integer column.
func (t *TableBuilder) Int(name string) *ColumnBuilder {
	colType := "INT"
	if t.driver == "postgres" {
		colType = "INTEGER"
	}
	c := &ColumnBuilder{name: name, typeName: colType}
	t.columns = append(t.columns, c)
	return c
}

// String adds a varchar column. Default length is 255.
func (t *TableBuilder) String(name string, length ...int) *ColumnBuilder {
	size := 255
	if len(length) > 0 && length[0] > 0 {
		size = length[0]
	}
	c := &ColumnBuilder{name: name, typeName: fmt.Sprintf("VARCHAR(%d)", size)}
	t.columns = append(t.columns, c)
	return c
}

// Timestamp adds a timestamp column.
func (t *TableBuilder) Timestamp(name string) *ColumnBuilder {
	c := &ColumnBuilder{name: name, typeName: "TIMESTAMP"}
	t.columns = append(t.columns, c)
	return c
}

// Timestamps adds created_at and updated_at columns.
func (t *TableBuilder) Timestamps() {
	t.Timestamp("created_at").DefaultRaw("CURRENT_TIMESTAMP")
	if t.driver == "mysql" {
		t.Timestamp("updated_at").DefaultRaw("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
		return
	}
	t.Timestamp("updated_at").DefaultRaw("CURRENT_TIMESTAMP")
}

// BuildCreateSQL renders CREATE TABLE IF NOT EXISTS SQL.
func (t *TableBuilder) BuildCreateSQL() string {
	parts := make([]string, 0, len(t.columns))
	for _, c := range t.columns {
		parts = append(parts, c.sql(t.driver))
	}
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)", t.table, strings.Join(parts, ",\n  "))
}

// NotNull marks the column as NOT NULL.
func (c *ColumnBuilder) NotNull() *ColumnBuilder {
	c.nullable = false
	return c
}

// Nullable marks the column as NULL.
func (c *ColumnBuilder) Nullable() *ColumnBuilder {
	c.nullable = true
	return c
}

// Default sets a quoted string default value.
func (c *ColumnBuilder) Default(value string) *ColumnBuilder {
	c.defaultSet = true
	c.defaultValue = value
	c.defaultIsRaw = false
	return c
}

// DefaultRaw sets a raw SQL default value/expression.
func (c *ColumnBuilder) DefaultRaw(value string) *ColumnBuilder {
	c.defaultSet = true
	c.defaultValue = value
	c.defaultIsRaw = true
	return c
}

// AutoIncrement marks the column as auto increment.
func (c *ColumnBuilder) AutoIncrement() *ColumnBuilder {
	c.autoIncrement = true
	return c
}

// PrimaryKey marks the column as primary key.
func (c *ColumnBuilder) PrimaryKey() *ColumnBuilder {
	c.primaryKey = true
	c.nullable = false
	return c
}

// Unique marks the column as unique.
func (c *ColumnBuilder) Unique() *ColumnBuilder {
	c.unique = true
	return c
}

func (c *ColumnBuilder) sql(driver string) string {
	typeName := c.typeName
	if driver == "postgres" && c.autoIncrement {
		typeName = "SERIAL"
	}

	pieces := []string{c.name, typeName}
	if driver == "mysql" && c.autoIncrement {
		pieces = append(pieces, "AUTO_INCREMENT")
	}
	if !c.nullable {
		pieces = append(pieces, "NOT NULL")
	}
	if c.defaultSet {
		if c.defaultIsRaw {
			pieces = append(pieces, "DEFAULT "+c.defaultValue)
		} else {
			pieces = append(pieces, fmt.Sprintf("DEFAULT '%s'", c.defaultValue))
		}
	}
	if c.unique {
		pieces = append(pieces, "UNIQUE")
	}
	if c.primaryKey {
		pieces = append(pieces, "PRIMARY KEY")
	}
	return strings.Join(pieces, " ")
}
