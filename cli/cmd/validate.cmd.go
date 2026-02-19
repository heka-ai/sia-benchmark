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

			excludedSections, err := cmd.Flags().GetStringArray("exclude-section")
			if err != nil {
				logger.Error().Err(err).Msg("Error getting excluded sections")
				return
			}

			ValidateExec(vllmModel, benchmarkModel, excludedSections)
		},
	}

	// Default to true so commands are printed during validate unless explicitly disabled
	cmd.Flags().Bool("vllm-command", false, "The model to use for the VLLM command")
	cmd.Flags().Bool("benchmark-command", false, "The model to use for the benchmark command")
	cmd.Flags().StringArray("exclude-section", []string{}, "Exclude a section from validation (can be used multiple times). Examples: --exclude-section benchmark --exclude-section vllm")

	return cmd
}

func formatCmd(binary string, args []string) string {
	if len(args) == 0 {
		return binary
	}
	return binary + " " + strings.Join(args, " ")
}

// This only validate that the TOML config file is valid
func ValidateExec(vllmModel bool, benchmarkModel bool, excludedSections []string) {
	logger.Info().Msg("Validating the config file")
	if len(excludedSections) > 0 {
		logger.Info().Strs("excluded_sections", excludedSections).Msg("Excluding sections from validation")
	}
	config.InitWithExclusions(excludedSections)

	if vllmModel {
		cfg := config.GetConfig()
		if cfg.VLLMConfig == nil {
			logger.Error().Msg("Cannot generate VLLM command: vllm section is excluded")
			return
		}

		localArgs, err := config.GenerateVLLMCommand(cfg.VLLMConfig)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating the VLLM command")
			return
		}

		token := ""
		if cfg.BenchmarkConfig != nil {
			token = cfg.BenchmarkConfig.Token
		}
		if token == "" {
			token = "${HF_TOKEN}"
		}
		logger.Info().Msg("VLLM command generated for your config:\n" + formatCmd("HF_TOKEN="+token+" vllm", localArgs))
	}

	if benchmarkModel {
		cfg := config.GetConfig()
		if cfg.BenchmarkConfig == nil {
			logger.Error().Msg("Cannot generate benchmark command: benchmark section is excluded")
			return
		}

		// Use EnginePort from config for validation output instead of a hardcoded default
		localArgs, err := config.GenerateBenchmarkCommand(&cfg, "127.0.0.1", cfg.BenchmarkConfig.EnginePort)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating the benchmark command")
			return
		}
		logger.Info().Msg("Benchmark command generated for your config:\n" + formatCmd("uv run python", localArgs))
	}
}
