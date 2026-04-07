package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tss182/migration/internal/migrator"
)

var seederRollbackStep int

var seedCmd = &cobra.Command{
	Use:   "seeder",
	Short: "Run all registered database seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.RunSeeders(conn)
	},
}

var seederRollbackCmd = &cobra.Command{
	Use:   "seeder:rollback",
	Short: "Rollback the last batch of database seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.RollbackSeeders(conn, seederRollbackStep)
	},
}

var seederStatusCmd = &cobra.Command{
	Use:   "seeder:status",
	Short: "Show the status of each seeder",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := connectDB()
		if err != nil {
			return err
		}
		defer conn.Close()
		return migrator.SeederStatus(conn)
	},
}

var seedAliasCmd = &cobra.Command{
	Use:    "seed",
	Short:  "Run all registered database seeders",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return seedCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(seederRollbackCmd)
	rootCmd.AddCommand(seederStatusCmd)
	rootCmd.AddCommand(seedAliasCmd)
	seederRollbackCmd.Flags().IntVar(&seederRollbackStep, "step", 1, "Number of batches to rollback")
}
