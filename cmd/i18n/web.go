package main

import (
	"log/slog"

	"github.com/VineLink-Lab/i18n/internal/web"
	"github.com/spf13/cobra"
)

var webCmd = cobra.Command{
	Use:   "web",
	Short: "Start web server for i18n management",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := web.Web(directory, port)
		if err != nil {
			slog.Error("code generation failed", slog.String("error", err.Error()))
			return err
		}
		return nil
	},
}
