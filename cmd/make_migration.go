package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/stubs"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		migName := timestamp + "_" + name
		structName := toPascalCase(name)
		filename := migName + ".go"

		stub, err := stubs.FS.ReadFile("migration.stub")
		if err != nil {
			return fmt.Errorf("failed to read migration stub: %w", err)
		}

		content := strings.ReplaceAll(string(stub), "{{NAME}}", migName)
		content = strings.ReplaceAll(content, "{{STRUCT}}", structName)

		dir := filepath.Join("database", "migrations")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		path := filepath.Join(dir, filename)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Created migration: %s\n", path)
		return nil
	},
}

func init() {
	makeCmd.AddCommand(makeMigrationCmd)
}

// toPascalCase converts snake_case or kebab-case to PascalCase.
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var sb strings.Builder
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		runes := []rune(part)
		sb.WriteRune(unicode.ToUpper(runes[0]))
		sb.WriteString(string(runes[1:]))
	}
	return sb.String()
}
