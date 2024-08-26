package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "timetastic",
	Short:   "Manage timetastic bookings from the command line",
	Long:    `Create your timetastic bookings from the command line`,
	Version: "development",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		rootCmd.Version = "unknown"
	}
	if info.Main.Version != "" {
		rootCmd.Version = info.Main.Version
	}
}
