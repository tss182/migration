package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
	Use:    "migrate",
	Short:  "Run all pending database migrations",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
