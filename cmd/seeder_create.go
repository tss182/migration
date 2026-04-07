/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/app"
)

// seederCreateCmd represents the seederCreate command
var seederCreateCmd = &cobra.Command{
	Use:    "seeder:create [name]",
	Short:  "Create a new seeder",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("seeder name is required")
		}
		app.CreateSeeder(".", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(seederCreateCmd)
}
