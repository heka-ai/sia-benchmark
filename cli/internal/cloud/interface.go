package cloud

import (
	"github.com/heka-ai/benchmark-cli/internal/config"
)

// Cloud is the interface for the different cloud providers
type Cloud interface {
	NewClient(config *config.Config) Cloud

	// Validate the cloud credentials
	ValidateCredentials() error

	// Create the two instances needed for the benchmark
	Create() error

	// Destroy the instances
	Destroy() error

	// Get the IP address of the LLM instance
	GetLLMInstanceIP() (string, error)

	// Get the IP address of the CPU instance
	GetBenchInstanceIP() (string, error)
}
