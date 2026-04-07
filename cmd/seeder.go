/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seederCmd represents the seeder command
var seederCmd = &cobra.Command{
	Use:    "seeder",
	Short:  "Run seeder files",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("seeder called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(seederCmd)
}
