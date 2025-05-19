package main

import (
	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/internal/config"
	"github.com/spf13/cobra"
)

// sends the command to run the benchmark
func BenchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the benchmark",
	}
}

func RunExec() {
	config.Init()
	c := config.GetConfig()

	client := bench.NewClient(c.APIKey)
	cloud := cloud_generator.NewCloud(&c)

	benchInstanceIP, err := cloud.GetBenchInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the bench instance IP")
	}

	err = client.RunBenchmark(benchInstanceIP)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot run benchmark on bench instance")
	}

	// todo: add a flag to wait for the benchmark to finish
	// todo: show the progress

	logger.Info().Msg("Benchmark started on the bench instance")
}
