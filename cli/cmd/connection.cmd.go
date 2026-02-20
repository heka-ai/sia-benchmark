package main

import (
	"github.com/spf13/cobra"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
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

// test the connection to the created instance(s); only checks instances that exist
func connect() {
	logger.Info().Msg("Trying to connect to the instances")
	config.InitInfra()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	client := bench.NewClient(c.GeneralConfig.APIKey)

	var llmIP, benchIP string
	if ip, err := cloud.GetLLMInstanceIP(); err == nil {
		llmIP = ip
		logger.Info().Str("ip", llmIP).Msg("LLM instance IP")
	} else {
		logger.Debug().Err(err).Msg("No LLM instance found")
	}
	if ip, err := cloud.GetBenchInstanceIP(); err == nil {
		benchIP = ip
		logger.Info().Str("ip", benchIP).Msg("Bench instance IP")
	} else {
		logger.Debug().Err(err).Msg("No bench instance found")
	}

	if llmIP == "" && benchIP == "" {
		logger.Fatal().Msg("No instances found; create at least one with 'bench create'")
	}

	if llmIP != "" {
		if err := client.HealthCheck(llmIP); err != nil {
			logger.Fatal().Err(err).Str("ip", llmIP).Msg("Cannot connect to the LLM instance")
		}
		logger.Info().Str("ip", llmIP).Msg("LLM instance connected")
	}
	if benchIP != "" {
		if err := client.HealthCheck(benchIP); err != nil {
			logger.Fatal().Err(err).Str("ip", benchIP).Msg("Cannot connect to the bench instance")
		}
		logger.Info().Str("ip", benchIP).Msg("Bench instance connected")
	}
}
