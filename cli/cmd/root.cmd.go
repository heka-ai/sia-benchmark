package main

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bench",
		Short: "Run business oriented LLM benchmarks",
		Long:  `Sia Benchmark is the cli to run business oriented LLM benchmarks`,
	}

	rootCmd.AddCommand(ValidateCmd())
	rootCmd.AddCommand(CredsCmd())
	rootCmd.AddCommand(InstanceCmd())
	rootCmd.AddCommand(ConnectionCmd())
	rootCmd.AddCommand(DeployCmd())
	rootCmd.AddCommand(BenchCmd())
	rootCmd.AddCommand(ResultsCmd())
	rootCmd.AddCommand(DestroyCmd())

	rootCmd.Flags().StringP("config", "c", "bench.toml", "Path to the config file")

	return rootCmd
}
