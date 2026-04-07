/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run migration commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migration called")
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}
