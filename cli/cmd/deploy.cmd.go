package main

import (
	"time"

	"github.com/spf13/cobra"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

// Deploy the model on the instance
func DeployCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "deploy",
		Short: "Run the setup of the model on the instance. This is long",
		Run: func(cmd *cobra.Command, args []string) {
			wait, err := cmd.Flags().GetBool("wait")
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get wait flag")
				return
			}

			deploy(wait)
		},
	}

	command.Flags().BoolP("wait", "w", false, "Wait for the LLM to be ready")

	return command
}

func deploy(wait bool) {
	logger.Info().Msg("Deploying and starting the LLM on the GPU instance")

	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	// Use port from config instead of flag
	port := c.BenchmarkConfig.EnginePort

	llmClient := bench.NewClient(c.GeneralConfig.APIKey)
	err = llmClient.Deploy(llmInstanceIP, c.GeneralConfig.InferenceEngine, port)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to deploy the LLM instance")
	}

	if wait {
		logger.Info().Msg("Waiting for the LLM to be ready")
		llmClient.WaitForLLM(llmInstanceIP, port)
		logger.Info().Msg("LLM is running")

		// Show logs while waiting
		logger.Info().Msg("Streaming LLM logs...")
		llmClient.FollowBenchLogs(llmInstanceIP, "llm", 1*time.Second, 5*time.Minute, true)
	} else {
		logger.Info().Msg("Model is downloading and being initialized on the LLM instance, wait a few minutes")
	}
}
