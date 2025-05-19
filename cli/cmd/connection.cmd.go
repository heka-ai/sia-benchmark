package main

import (
	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/internal/config"
	"github.com/spf13/cobra"
)

func ConnectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connection",
		Short: "Validate the connection to the created instance",
		Run: func(cmd *cobra.Command, args []string) {
			connect()
		},
	}
}

// test the connection to the two instances
func connect() {
	logger.Info().Msg("Trying to connect to the instances")
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	client := bench.NewClient(c.APIKey)

	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	benchInstanceIP, err := cloud.GetBenchInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the bench instance IP")
	}

	err = client.HealthCheck(llmInstanceIP)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot connect to the LLM instance")
	}
	logger.Info().Str("ip", llmInstanceIP).Msg("LLM instance connected")

	err = client.HealthCheck(benchInstanceIP)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot connect to the bench instance")
	}

	logger.Info().Str("ip", benchInstanceIP).Msg("Bench instance connected")
}
