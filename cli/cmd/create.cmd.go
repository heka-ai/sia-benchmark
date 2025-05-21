package main

import (
	"github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

func InstanceCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create the instance to run the benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			wait, err := cmd.Flags().GetBool("wait")
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get wait flag")
				return
			}

			create(wait)
		},
	}

	command.Flags().BoolP("wait", "w", false, "Wait for the instances to be ready")

	return command
}

func create(wait bool) {
	logger.Info().Msg("Creating the instances to run the benchmark")
	config.Init()

	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	cloud.Create()

	if wait {
		logger.Info().Msg("Waiting for the instances to be ready")
		bench := bench.NewClient(c.APIKey)

		benchIP, err := cloud.GetBenchInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to get the bench instance IP")
		}

		llmIP, err := cloud.GetLLMInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to get the LLM instance IP")
		}

		bench.WaitForInstances(benchIP, llmIP)

		logger.Info().Msg("Instances are ready")
	}
}
