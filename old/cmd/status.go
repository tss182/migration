package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/migrator"
)

var statusCmd = &cobra.Command{
	Use:   "migration:status",
	Short: "Show the status of each migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.Status(conn)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
