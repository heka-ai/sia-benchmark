package main

import (
	"github.com/spf13/cobra"
)

// Output the results of the benchmark
func ResultsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "results",
		Short: "Get the results of the benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			results()
		},
	}
}

func results() {
	logger.Fatal().Msg("Not implemented")
}
