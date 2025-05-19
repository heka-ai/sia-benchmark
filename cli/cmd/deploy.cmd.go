package main

import (
	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

// Deploy the model on the instance
func DeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Run the setup of the model on the instance. This is long",
		Run: func(cmd *cobra.Command, args []string) {
			deploy()
		},
	}
}

func deploy() {
	logger.Info().Msg("Deploying and starting the LLM on the GPU instance")

	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	llmClient := bench.NewClient(c.APIKey)
	err = llmClient.Deploy(llmInstanceIP, c.InferenceEngine)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to deploy the LLM instance")
	}

	// todo: add a flag to wait for the llm to be ready
	// todo: add a command to get the logs from the LLM instance / engine

	logger.Info().Msg("Model is downloading and being initialized on the LLM instance, wait a few minutes")
}
