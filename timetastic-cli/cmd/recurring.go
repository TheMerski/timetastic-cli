package cmd

import (
	"github.com/spf13/cobra"
	"github.com/themerski/timetastic-cli/internal/api"
	"github.com/themerski/timetastic-cli/internal/flows"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "recurring",
	Short: "Create a recurring booking",
	Long:  `Book weekly or bi-weekly recurring leaves`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewTimetasticClient()

		// Only run the flow if the client is not nil
		// Errors are already logged in the NewTimetasticClient function
		if client != nil {
			flows.BookRecurringLeave(client)
		}
	},
}
