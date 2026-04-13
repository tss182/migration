package app

import "database/sql"

var (
	Driver = "mysql" // default driver, can be overridden by environment variable
	Tables []string
	Db     *sql.DB
)

func CreateMigration(baseDir string, name string) (string, error) {
	return CreateFile(baseDir, name, "migration")
}

func Init(name string) {
	Tables = append(Tables, name)
}
