package main

import (
	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

func SetupInstanceCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "setup-instance",
		Short: "Install dependencies on a created instance (LLM or benchmark)",
		Run: func(cmd *cobra.Command, args []string) {
			llm, _ := cmd.Flags().GetBool("llm")
			benchmark, _ := cmd.Flags().GetBool("benchmark")

			if !llm && !benchmark {
				logger.Fatal().Msg("Specify --llm or --benchmark")
			}
			if llm && benchmark {
				logger.Fatal().Msg("Specify only one of --llm or --benchmark")
			}

			if llm {
				setupInstance("llm")
			} else {
				setupInstance("benchmark")
			}
		},
	}

	command.Flags().Bool("llm", false, "Setup the LLM (GPU) instance")
	command.Flags().Bool("benchmark", false, "Setup the benchmark (CPU) instance")

	return command
}

func setupInstance(setupType string) {
	config.InitInfra()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	client := bench.NewClient(c.GeneralConfig.APIKey)

	var ip string
	var err error

	if setupType == "llm" {
		ip, err = cloud.GetLLMInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
		}
	} else {
		ip, err = cloud.GetBenchInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot get the bench instance IP")
		}
	}

	logger.Info().Str("ip", ip).Str("type", setupType).Msg("Setting up instance")

	err = client.SetupInstance(ip, setupType)
	if err != nil {
		logger.Fatal().Err(err).Msg("Setup failed")
	}

	logger.Info().Str("ip", ip).Str("type", setupType).Msg("Instance setup complete")
}
