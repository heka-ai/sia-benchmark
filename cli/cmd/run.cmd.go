package main

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	bench "github.com/heka-ai/benchmark-cli/internal/bench"
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

// sends the command to run the benchmark
func BenchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			RunExec()
		},
	}

	cmd.Flags().Int("vllm-port", 8000, "Port where vLLM listens (default 8000)")
	cmd.Flags().String("result-filename", "./metrics.json", "Where to write benchmark results")
	cmd.Flags().Bool("wait", true, "Stream benchmark logs until completion")

	// Bind flags to Viper keys for global access
	_ = viper.BindPFlag("vllm_port", cmd.Flags().Lookup("vllm-port"))
	_ = viper.BindPFlag("result_filename", cmd.Flags().Lookup("result-filename"))
	_ = viper.BindPFlag("wait", cmd.Flags().Lookup("wait"))

	// Ensure env like VLLM_PORT works (optional)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return cmd
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

	llmInstanceIP, err := cloud.GetLLMInstanceIP()
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot get the LLM instance IP")
	}

	// Read values via Viper (flags/env/config precedence)
	port := viper.GetInt("vllm_port")
	result := viper.GetString("result_filename")
	if strings.TrimSpace(result) == "" {
		result = "./metrics.json"
	}
	wait := viper.GetBool("wait")

	err = client.RunBenchmark(benchInstanceIP, llmInstanceIP, c.InferenceEngine, port, result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot run benchmark on bench instance")
	}

	logger.Info().Msg("Benchmark started on the bench instance")
	if wait {
		_ = client.FollowBenchLogsWithTimeout(benchInstanceIP, "bench", 1*time.Second, 5*time.Minute, true)
	}
}
