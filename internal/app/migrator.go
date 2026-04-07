package app

var Driver = "mysql" // default driver, can be overridden by environment variable

func CreateMigration(baseDir string, name string) (string, error) {
	return CreateFile(baseDir, name, "migration")
}
