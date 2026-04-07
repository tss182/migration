/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/app"
)

// migrationCreateCmd represents the migrationCreate command
var migrationCreateCmd = &cobra.Command{
	Use:    "migration:create [name]",
	Short:  "Create a new migration",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("migration name is required")
		}
		app.CreateMigration(".", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrationCreateCmd)
}
