package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	directory string
	port      string
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

func init() {
	rootCmd.Flags().StringVarP(&directory, "dir", "d", ".", "Directory containing input files")
	err := rootCmd.MarkFlagRequired("dir")
	if err != nil {
		slog.Error("failed to mark 'dir' flag as required", slog.String("error", err.Error()))
		return
	}
	generateCmd.Flags().StringVarP(&directory, "dir", "d", ".", "Directory containing input files")
	err = generateCmd.MarkFlagRequired("dir")
	if err != nil {
		slog.Error("failed to mark 'dir' flag as required", slog.String("error", err.Error()))
		return
	}
	webCmd.Flags().StringVarP(&directory, "dir", "d", ".", "Directory containing input files")
	err = webCmd.MarkFlagRequired("dir")
	if err != nil {
		slog.Error("failed to mark 'dir' flag as required", slog.String("error", err.Error()))
		return
	}
	webCmd.Flags().StringVarP(&port, "port", "p", "1180", "Port for the web server")
	rootCmd.AddCommand(&generateCmd, &webCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("execution failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
