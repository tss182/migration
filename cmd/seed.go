package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/config"
	"github.com/tss182/migration/internal/db"
	"github.com/tss182/migration/internal/migrator"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run all registered database seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		conn, err := db.Connect(cfg)
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.RunSeeders(conn)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
