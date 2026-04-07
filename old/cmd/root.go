package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/config"
	"github.com/tss182/migration/internal/db"
	"github.com/tss182/migration/internal/migrator"
)

var rootCmd = &cobra.Command{
	Use:          "migration",
	Short:        "migration and seeder CLI tool for Go",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.Run(conn)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		// Load .env file if present; ignore error when file doesn't exist
		_ = godotenv.Load()
	})
}

func connectDB() (*sql.DB, error) {
	cfg := config.Load()
	conn, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	migrator.Driver = cfg.Driver
	return conn, nil
}
