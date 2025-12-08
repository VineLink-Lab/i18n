package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			slog.Error("failed to display help", slog.String("error", err.Error()))
		}
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("execution failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
