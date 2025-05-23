package main

import (
	"fmt"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

func LogCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "logs",
		Short: "Get the logs of the benchmark",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logsType := args[0]
			logs(logsType)
		},
	}

	command.Flags().BoolP("follow", "f", false, "Follow the logs")
	command.Flags().StringP("type", "t", "llm", "The type of logs to get")

	return command
}

func logs(logsType string) {
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	llmClient := bench.NewClient(c.APIKey)
	logs, err := llmClient.GetLogs(llmInstanceIP, logsType)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the logs")
	}

	fmt.Println(logs)
}
