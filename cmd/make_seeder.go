package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/stubs"
)

var makeSeederCmd = &cobra.Command{
	Use:   "seeder [name]",
	Short: "Create a new seeder file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		structName := toPascalCase(name)
		filename := time.Now().Format("20060102150405") + "_" + strings.ToLower(name) + ".go"

		stub, err := stubs.FS.ReadFile("seeder.stub")
		if err != nil {
			return fmt.Errorf("failed to read seeder stub: %w", err)
		}

		content := strings.ReplaceAll(string(stub), "{{NAME}}", structName)
		content = strings.ReplaceAll(content, "{{STRUCT}}", structName)

		dir := filepath.Join("database", "seeders")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		path := filepath.Join(dir, filename)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Created seeder: %s\n", path)
		return nil
	},
}

func init() {
	makeCmd.AddCommand(makeSeederCmd)
}
