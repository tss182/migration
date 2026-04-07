package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/config"
	"github.com/tss182/migration/internal/db"
	"github.com/tss182/migration/internal/migrator"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run all pending database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		conn, err := db.Connect(cfg)
		if err != nil {
			return err
		}
		defer conn.Close()
		migrator.Driver = cfg.Driver
		return migrator.Run(conn)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
