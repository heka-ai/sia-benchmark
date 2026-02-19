package main

import (
	"time"

	"github.com/spf13/cobra"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

func LogsCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "logs",
		Short: "Get logs from LLM or Benchmark instances",
		Long:  `Get logs from either the LLM instance or the Benchmark instance`,
		Example: `
		bench logs --llm
		bench logs --benchmark
		bench logs --llm --follow
		bench logs --benchmark --follow --interval 2s
		`,
		Run: func(cmd *cobra.Command, args []string) {
			llm, err := cmd.Flags().GetBool("llm")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get llm flag")
			}

			benchmark, err := cmd.Flags().GetBool("benchmark")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get benchmark flag")
			}

			follow, err := cmd.Flags().GetBool("follow")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get follow flag")
			}

			intervalStr, err := cmd.Flags().GetString("interval")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get interval flag")
			}

			interval, err := time.ParseDuration(intervalStr)
			if err != nil {
				logger.Fatal().Err(err).Msg("Invalid interval format. Use formats like '2s', '500ms', '1m'")
			}

			if !llm && !benchmark {
				logger.Fatal().Msg("You must specify either --llm or --benchmark")
			}

			if llm && benchmark {
				logger.Fatal().Msg("You cannot specify both --llm and --benchmark at the same time")
			}

			logs(llm, benchmark, follow, interval)
		},
	}

	command.Flags().Bool("llm", false, "Get logs from the LLM instance")
	command.Flags().Bool("benchmark", false, "Get logs from the Benchmark instance")
	command.Flags().BoolP("follow", "f", false, "Follow the logs in real-time")
	command.Flags().StringP("interval", "i", "1s", "Interval between log checks when following (e.g., '2s', '500ms', '1m')")

	return command
}

func logs(isLLM, isBenchmark, follow bool, interval time.Duration) {
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	client := bench.NewClient(c.GeneralConfig.APIKey)

	var instanceIP string
	var instanceType string
	var err error

	if isLLM {
		instanceIP, err = cloud.GetLLMInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
		}
		instanceType = "llm"
		logger.Info().Str("ip", instanceIP).Msg("Getting LLM instance logs")
	} else {
		instanceIP, err = cloud.GetBenchInstanceIP()
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot get the Benchmark instance IP")
		}
		instanceType = "bench"
		logger.Info().Str("ip", instanceIP).Msg("Getting Benchmark instance logs")
	}

	if follow {
		logger.Info().Str("type", instanceType).Msg("Following logs in real-time")
		err = client.FollowBenchLogs(instanceIP, instanceType, interval, 5*time.Minute, true)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to follow logs")
		}
	} else {
		logs, err := client.GetLogs(instanceIP, instanceType)
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot get the logs")
		}
		logger.Info().Msg(logs)
	}
}
