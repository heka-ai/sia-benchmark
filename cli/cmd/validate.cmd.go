package main

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/heka-ai/benchmark-cli/pkg/config"
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

	// Default to true so commands are printed during validate unless explicitly disabled
	cmd.Flags().Bool("vllm-command", false, "The model to use for the VLLM command")
	cmd.Flags().Bool("benchmark-command", false, "The model to use for the benchmark command")

	return cmd
}

func formatCmd(binary string, args []string) string {
	if len(args) == 0 {
		return binary
	}
	return binary + " " + strings.Join(args, " ")
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

		logger.Info().Msg("VLLM command generated for your config:\n" + formatCmd("HF_TOKEN="+cfg.BenchmarkConfig.Token+" vllm", localArgs))
	}

	if benchmarkModel {
		cfg := config.GetConfig()

		// Use EnginePort from config for validation output instead of a hardcoded default
		localArgs, err := config.GenerateBenchmarkCommand(&cfg, "127.0.0.1", cfg.BenchmarkConfig.EnginePort)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating the benchmark command")
			return
		}
		logger.Info().Msg("Benchmark command generated for your config:\n" + formatCmd("/opt/pytorch/bin/python3", localArgs))
	}
}
