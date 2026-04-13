package builder

import (
	"fmt"

	"github.com/tss182/migration/internal/app"
)

type (
	Table struct {
		driver  string
		name    string
		columns []*Column
	}
	Column struct {
		name          string
		typeName      string
		nullable      bool
		defaultSet    bool
		defaultValue  string
		autoIncrement bool
		primaryKey    bool
		unique        bool
	}
)

func Create(tableName string, build func(t *Table)) error {
	tb := NewTable(app.Driver, tableName)
	build(tb)
	_, err := app.Db.Exec(tb.build())
	return err
}

func Drop(table string) error {
	_, err := app.Db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	return err
}

func NewTable(driver, table string) *Table {
	return &Table{driver: driver, name: table}
}

func (t *Table) build() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", t.name, "")
}

func (t *Table) Int(name string) *Column {
	colType := "INT"
	if t.driver == "postgres" {
		colType = "INTEGER"
	}
	c := &Column{name: name, typeName: colType}
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) String(name string, length int) *Column {
	c := &Column{name: name, typeName: fmt.Sprintf("VARCHAR(%d)", length)}
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) Timestamp(name string) *Column {
	var c *Column
	if t.driver == "postgres" {
		c = &Column{name: name, typeName: "TIMESTAMPTZ"}
	} else {
		c = &Column{name: name, typeName: "TIMESTAMP"}
	}
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) TimestampAuto() {
	t.Timestamp("created_at").NotNull()
	t.Timestamp("updated_at").NotNull()
	t.Timestamp("deleted_at")
}

func (c *Column) NotNull() *Column {
	c.nullable = false
	return c
}

func (c *Column) Default(value string) *Column {
	c.defaultSet = true
	c.defaultValue = value
	return c
}

func (c *Column) AutoIncrement() *Column {
	c.autoIncrement = true
	return c
}

func (c *Column) PrimaryKey() *Column {
	c.primaryKey = true
	return c
}

func (c *Column) Unique() *Column {
	c.unique = true
	return c
}
