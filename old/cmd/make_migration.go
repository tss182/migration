package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/generator"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "migration:create [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := generator.CreateMigration(".", args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Created migration: %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
	makeCmd.AddCommand(&cobra.Command{
		Use:    "migration [name]",
		Short:  "Create a new migration file",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeMigrationCmd.RunE(cmd, args)
		},
	})
}
