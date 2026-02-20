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
			wait, _ := cmd.Flags().GetBool("wait")
			llm, _ := cmd.Flags().GetBool("llm")
			bench, _ := cmd.Flags().GetBool("benchmark")

			// No flags → both. --llm only → LLM. --benchmark only → benchmark. Both flags → both.
			createLLM := llm || !bench
			createBench := bench || !llm

			create(createLLM, createBench, wait)
		},
	}

	command.Flags().BoolP("wait", "w", false, "Wait for the instances to be ready")
	command.Flags().Bool("llm", false, "Create only the LLM (GPU) instance")
	command.Flags().Bool("benchmark", false, "Create only the benchmark (CPU) instance")

	return command
}

func create(createLLM, createBench, wait bool) {
	logger.Info().Bool("llm", createLLM).Bool("benchmark", createBench).Msg("Creating instances")
	config.InitInfra()

	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	if err := cloud.Create(createLLM, createBench); err != nil {
		logger.Fatal().Err(err).Msg("Create failed")
	}

	if wait {
		logger.Info().Msg("Waiting for the instances to be ready")
		client := bench.NewClient(c.GeneralConfig.APIKey)

		var benchIP, llmIP string
		if createBench {
			var err error
			benchIP, err = cloud.GetBenchInstanceIP()
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get the bench instance IP")
			}
		}
		if createLLM {
			var err error
			llmIP, err = cloud.GetLLMInstanceIP()
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to get the LLM instance IP")
			}
		}

		client.WaitForInstances(benchIP, llmIP)
		logger.Info().Msg("Instances are ready")
	}
}
