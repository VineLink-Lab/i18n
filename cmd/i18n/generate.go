package main

import (
	"log/slog"

	"github.com/VineLink-Lab/i18n/internal/generate"

	"github.com/spf13/cobra"
)

var (
	directory string
)

var generateCmd = cobra.Command{
	Use:   "generate",
	Short: "Generate go code from input Directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := generate.Generate(directory)
		if err != nil {
			slog.Error("code generation failed", slog.String("error", err.Error()))
			return err
		}
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&directory, "dir", "d", ".", "Directory containing input files")
	err := generateCmd.MarkFlagRequired("dir")
	if err != nil {
		slog.Error("failed to mark 'dir' flag as required", slog.String("error", err.Error()))
		return
	}
	rootCmd.AddCommand(&generateCmd)
}
