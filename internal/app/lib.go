package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/tss182/migration/internal/stubs"
)

var (
	MigrationsDir = "database/migrations"
	SeedersDir    = "database/seeders"
)

func CreateFile(baseDir , name, fileType string) (string, error) {
	var folderPath string
	switch fileType {
	case "migration":
		folderPath = filepath.Join(baseDir, MigrationsDir)
	case "seeder":
		folderPath = filepath.Join(baseDir, SeedersDir)
	default:
		return "", fmt.Errorf("unknown file type: %s", fileType)
	}

	if err := createFolderIFNotExist(folderPath); err != nil {
		return "", err
	}
	timestamp := time.Now().UnixMicro()
	fileName := fmt.Sprintf("%d_%s", timestamp, toSnakeCase(name))
	structName := toPascalCase(name)
	filename := fileName + ".go"

	stub, err := stubs.FS.ReadFile(fileType + ".stub")
	if err != nil {
		return "", fmt.Errorf("failed to read %s stub: %w", fileType, err)
	}

	content := strings.ReplaceAll(string(stub), "{{NAME}}", fileName)
	content = strings.ReplaceAll(content, "{{STRUCT}}", structName)

	path := filepath.Join(folderPath, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return path, nil
}

func createFolderIFNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

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

func toSnakeCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var sb strings.Builder
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		if i > 0 {
			sb.WriteRune('_')
		}
		runes := []rune(part)
		sb.WriteRune(unicode.ToLower(runes[0]))
		sb.WriteString(string(runes[1:]))
	}
	return sb.String()
}
