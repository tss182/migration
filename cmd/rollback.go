package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/config"
	"github.com/tss182/migration/internal/db"
	"github.com/tss182/migration/internal/migrator"
)

var rollbackStep int

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last batch of database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		conn, err := db.Connect(cfg)
		if err != nil {
			return err
		}
		defer conn.Close()
		migrator.Driver = cfg.Driver
		return migrator.Rollback(conn, rollbackStep)
	},
}

func init() {
	rollbackCmd.Flags().IntVar(&rollbackStep, "step", 1, "Number of batches to rollback")
	migrateCmd.AddCommand(rollbackCmd)
}
