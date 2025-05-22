package main

import (
	"encoding/json"
	"os"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

// Output the results of the benchmark
func ResultsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "results",
		Short: "Get the results of the benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				logger.Fatal().Msgf("Error getting file: %v", err)
			}

			if file == "" {
				logger.Fatal().Msg("file flag is required")
			}

			results(file)
		},
	}

	cmd.Flags().StringP("file", "f", "", "The file to write the results to")

	return cmd
}

func results(file string) {
	config.Init()
	config := config.GetConfig()
	client := bench.NewClient(config.APIKey)

	cloud := cloud_generator.NewCloud(&config)

	benchInstanceIP, err := cloud.GetBenchInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the bench instance IP")
	}

	results, err := client.GetResults(benchInstanceIP, config.InferenceEngine)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the results")
	}

	json, err := json.Marshal(results)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot marshal the results")
	}

	os.WriteFile(file, json, 0644)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot write the results to the file")
	}

	logger.Info().Msgf("Results written to %s", file)
}
