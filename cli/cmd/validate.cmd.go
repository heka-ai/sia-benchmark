package main

import (
	"strings"

	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

// validate the config of the benchmark making sure all the right credentials are set
func ValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the benchmark",
		Long:  `Validate the benchmark by checking the toml config file`,
		Example: `
		bench validate
		`,
		Run: func(cmd *cobra.Command, args []string) {
			vllmModel, err := cmd.Flags().GetBool("vllm-command")
			if err != nil {
				logger.Error().Err(err).Msg("Error getting the VLLM model")
				return
			}

			benchmarkModel, err := cmd.Flags().GetBool("benchmark-command")
			if err != nil {
				logger.Error().Err(err).Msg("Error getting the benchmark model")
				return
			}

			ValidateExec(vllmModel, benchmarkModel)
		},
	}

	cmd.Flags().Bool("vllm-command", false, "The model to use for the VLLM command")
	cmd.Flags().Bool("benchmark-command", false, "The model to use for the benchmark command")

	return cmd
}

// This only validate that the TOML config file is valid
func ValidateExec(vllmModel bool, benchmarkModel bool) {
	logger.Info().Msg("Validating the config file")
	config.Init()

	if vllmModel {
		cfg := config.GetConfig()

		localArgs, err := config.GenerateVLLMCommand(cfg.VLLMConfig)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating the VLLM command")
			return
		}

		logger.Info().Str("command", "vllm "+strings.Join(localArgs, " ")).Msg("VLLM command generated for your config")
	}

	if benchmarkModel {
		cfg := config.GetConfig()

		localArgs, err := config.GenerateBenchmarkCommand(&cfg, "127.0.0.1")
		if err != nil {
			logger.Error().Err(err).Msg("Error generating the benchmark command")
			return
		}

		logger.Info().Str("command", "/opt/pytorch/bin/python3 "+strings.Join(localArgs, " ")).Msg("Benchmark command generated for your config")
	}
}
