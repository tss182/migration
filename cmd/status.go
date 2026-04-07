package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/config"
	"github.com/tss182/migration/internal/db"
	"github.com/tss182/migration/internal/migrator"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of each migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		conn, err := db.Connect(cfg)
		if err != nil {
			return err
		}
		defer conn.Close()
		migrator.Driver = cfg.Driver
		return migrator.Status(conn)
	},
}

func init() {
	migrateCmd.AddCommand(statusCmd)
}
