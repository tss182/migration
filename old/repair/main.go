package main

import (
	"fmt"
	"os"

	"github.com/tss182/migration/internal/generator"
	"github.com/tss182/migration/internal/project"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || (len(args) == 1 && args[0] == "repair") {
		if err := project.EnsureRegistryPackages("."); err != nil {
			return err
		}
		fmt.Println("Registry packages are ready:")
		fmt.Println("- database/migrations/register.go")
		fmt.Println("- database/seeders/register.go")
		return nil
	}

	if len(args) == 3 && args[0] == "make" && args[1] == "migration" {
		path, err := generator.CreateMigration(".", args[2])
		if err != nil {
			return err
		}
		fmt.Printf("Created migration: %s\n", path)
		return nil
	}

	if len(args) == 3 && args[0] == "make" && args[1] == "seeder" {
		path, err := generator.CreateSeeder(".", args[2])
		if err != nil {
			return err
		}
		fmt.Printf("Created seeder: %s\n", path)
		return nil
	}

	return fmt.Errorf("usage:\n  go run ./tools/repair\n  go run ./tools/repair repair\n  go run ./tools/repair make migration <name>\n  go run ./tools/repair make seeder <name>")
}
