package main

import (
	"github.com/heka-ai/benchmark-cli/internal/config"
	"github.com/spf13/cobra"
)

// validate the config of the benchmark making sure all the right credentials are set
func ValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the benchmark",
		Long:  `Validate the benchmark by checking the toml config file`,
		Example: `
		bench validate
		`,
		Run: func(cmd *cobra.Command, args []string) {
			ValidateExec()
		},
	}
}

// This only validate that the TOML config file is valid
func ValidateExec() {
	logger.Info().Msg("Validating the config file")
	config.Init()
}
