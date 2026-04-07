package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/generator"
)

var makeSeederCmd = &cobra.Command{
	Use:   "seeder:create [name]",
	Short: "Create a new seeder file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := generator.CreateSeeder(".", args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Created seeder: %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeSeederCmd)
	makeCmd.AddCommand(&cobra.Command{
		Use:    "seeder [name]",
		Short:  "Create a new seeder file",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeSeederCmd.RunE(cmd, args)
		},
	})
}
