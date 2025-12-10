package main

import (
	"log/slog"

	"github.com/VineLink-Lab/i18n/internal/generate"

	"github.com/spf13/cobra"
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
