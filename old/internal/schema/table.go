package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

// TableBuilder provides a fluent API for building CREATE TABLE SQL.
type TableBuilder struct {
	driver      string
	table       string
	columns     []*ColumnBuilder
	indexes     []string
	uniqueKeys  []string
	constraints []string
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

// SoftDeletes adds a deleted_at column for soft delete support.
func (t *TableBuilder) SoftDeletes() {
	t.Timestamp("deleted_at").Nullable()
}

// BuildCreateSQL renders CREATE TABLE IF NOT EXISTS SQL.
func (t *TableBuilder) BuildCreateSQL() string {
	parts := make([]string, 0, len(t.columns)+len(t.uniqueKeys)+len(t.indexes)+len(t.constraints))
	for _, c := range t.columns {
		parts = append(parts, c.sql(t.driver))
	}
	parts = append(parts, t.uniqueKeys...)
	parts = append(parts, t.indexes...)
	parts = append(parts, t.constraints...)
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)", t.table, strings.Join(parts, ",\n  "))
}

// Index adds a regular index on one or more columns.
func (t *TableBuilder) Index(columns ...string) *TableBuilder {
	if len(columns) == 0 {
		return t
	}
	idxName := t.table + "_" + strings.Join(columns, "_") + "_idx"
	idxDef := fmt.Sprintf("INDEX %s (%s)", idxName, strings.Join(columns, ", "))
	t.indexes = append(t.indexes, idxDef)
	return t
}

// UniqueIndex adds a unique index on one or more columns.
func (t *TableBuilder) UniqueIndex(columns ...string) *TableBuilder {
	if len(columns) == 0 {
		return t
	}
	idxName := t.table + "_" + strings.Join(columns, "_") + "_unique"
	idxDef := fmt.Sprintf("UNIQUE INDEX %s (%s)", idxName, strings.Join(columns, ", "))
	t.uniqueKeys = append(t.uniqueKeys, idxDef)
	return t
}

// ForeignKey adds a foreign key constraint.
func (t *TableBuilder) ForeignKey(column, refTable, refColumn string) *TableBuilder {
	fkName := t.table + "_" + column + "_fk"
	fkDef := fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)", fkName, column, refTable, refColumn)
	t.constraints = append(t.constraints, fkDef)
	return t
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

// AlterBuilder provides a fluent API for ALTER TABLE statements.
type AlterBuilder struct {
	driver      string
	table       string
	alterations []string
}

// AlterTable initializes an alter table builder.
func AlterTable(db *sql.DB, driver, table string, build func(a *AlterBuilder)) error {
	ab := NewAlter(driver, table)
	build(ab)
	if len(ab.alterations) == 0 {
		return nil
	}
	sql := ab.BuildAlterSQL()
	_, err := db.Exec(sql)
	return err
}

// NewAlter initializes a fluent alter table builder.
func NewAlter(driver, table string) *AlterBuilder {
	return &AlterBuilder{driver: driver, table: table}
}

// AddColumn adds a new column to the table.
func (a *AlterBuilder) AddColumn(col *ColumnBuilder) *AlterBuilder {
	stmt := fmt.Sprintf("ADD COLUMN %s %s", col.name, col.sql(a.driver))
	a.alterations = append(a.alterations, stmt)
	return a
}

// DropColumn removes a column from the table.
func (a *AlterBuilder) DropColumn(name string) *AlterBuilder {
	stmt := fmt.Sprintf("DROP COLUMN %s", name)
	a.alterations = append(a.alterations, stmt)
	return a
}

// RenameColumn renames a column.
func (a *AlterBuilder) RenameColumn(oldName, newName string) *AlterBuilder {
	var stmt string
	if a.driver == "postgres" {
		stmt = fmt.Sprintf("RENAME COLUMN %s TO %s", oldName, newName)
	} else {
		// MySQL syntax: CHANGE old_name new_name type
		stmt = fmt.Sprintf("CHANGE COLUMN %s %s VARCHAR(255)", oldName, newName)
	}
	a.alterations = append(a.alterations, stmt)
	return a
}

// ModifyColumn changes a column definition.
func (a *AlterBuilder) ModifyColumn(col *ColumnBuilder) *AlterBuilder {
	var stmt string
	if a.driver == "postgres" {
		stmt = fmt.Sprintf("ALTER COLUMN %s TYPE %s", col.name, col.typeName)
	} else {
		stmt = fmt.Sprintf("MODIFY COLUMN %s %s", col.name, col.sql(a.driver))
	}
	a.alterations = append(a.alterations, stmt)
	return a
}

// AddIndex adds an index to the table.
func (a *AlterBuilder) AddIndex(columns ...string) *AlterBuilder {
	if len(columns) == 0 {
		return a
	}
	idxName := a.table + "_" + strings.Join(columns, "_") + "_idx"
	stmt := fmt.Sprintf("ADD INDEX %s (%s)", idxName, strings.Join(columns, ", "))
	a.alterations = append(a.alterations, stmt)
	return a
}

// DropIndex removes an index from the table.
func (a *AlterBuilder) DropIndex(indexName string) *AlterBuilder {
	var stmt string
	if a.driver == "postgres" {
		stmt = fmt.Sprintf("DROP INDEX %s", indexName)
	} else {
		stmt = fmt.Sprintf("DROP INDEX %s ON %s", indexName, a.table)
	}
	a.alterations = append(a.alterations, stmt)
	return a
}

// AddForeignKey adds a foreign key constraint.
func (a *AlterBuilder) AddForeignKey(column, refTable, refColumn string) *AlterBuilder {
	fkName := a.table + "_" + column + "_fk"
	stmt := fmt.Sprintf("ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)", fkName, column, refTable, refColumn)
	a.alterations = append(a.alterations, stmt)
	return a
}

// DropForeignKey removes a foreign key constraint.
func (a *AlterBuilder) DropForeignKey(constraintName string) *AlterBuilder {
	var stmt string
	if a.driver == "postgres" {
		stmt = fmt.Sprintf("DROP CONSTRAINT %s", constraintName)
	} else {
		stmt = fmt.Sprintf("DROP FOREIGN KEY %s", constraintName)
	}
	a.alterations = append(a.alterations, stmt)
	return a
}

// BuildAlterSQL renders ALTER TABLE SQL combining all alterations.
func (a *AlterBuilder) BuildAlterSQL() string {
	return fmt.Sprintf("ALTER TABLE %s %s", a.table, strings.Join(a.alterations, ", "))
}
