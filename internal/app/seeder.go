package app

func CreateSeeder(baseDir string, name string) (string, error) {
	return CreateFile(baseDir, name, "seeder")
}
