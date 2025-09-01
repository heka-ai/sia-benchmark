package main

import (
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

			port, err := cmd.Flags().GetInt("vllm-port")
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get vllm-port flag")
				return
			}

			deploy(wait, port)
		},
	}

	command.Flags().BoolP("wait", "w", false, "Wait for the LLM to be ready")
	command.Flags().Int("vllm-port", 8000, "Port where vLLM will listen (default 8000)")

	return command
}

func deploy(wait bool, port int) {
	logger.Info().Msg("Deploying and starting the LLM on the GPU instance")

	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	llmClient := bench.NewClient(c.APIKey)
	err = llmClient.Deploy(llmInstanceIP, c.InferenceEngine, port)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to deploy the LLM instance")
	}

	if wait {
		logger.Info().Msg("Waiting for the LLM to be ready")
		llmClient.WaitForLLM(llmInstanceIP, port)
		logger.Info().Msg("LLM is running")
	} else {
		logger.Info().Msg("Model is downloading and being initialized on the LLM instance, wait a few minutes")
	}

	// todo: add a command to get the logs from the LLM instance / engine

}
