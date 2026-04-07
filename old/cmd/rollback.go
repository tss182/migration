package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/migrator"
)

var rollbackStep int

var rollbackCmd = &cobra.Command{
	Use:   "migration:rollback",
	Short: "Rollback the last batch of database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.Rollback(conn, rollbackStep)
	},
}

func init() {
	rollbackCmd.Flags().IntVar(&rollbackStep, "step", 1, "Number of batches to rollback")
	rootCmd.AddCommand(rollbackCmd)
}
