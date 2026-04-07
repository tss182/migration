package cmd

import "github.com/spf13/cobra"

var makeCmd = &cobra.Command{
	Use:    "make",
	Short:  "Generate migration or seeder stub files",
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
