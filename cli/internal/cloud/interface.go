package cloud

import (
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

// Cloud is the interface for the different cloud providers
type Cloud interface {
	NewClient(config *config.Config) Cloud

	// Validate the cloud credentials
	ValidateCredentials() error

	// Create instances. createLLM and createBench control which to create; both true = both.
	Create(createLLM, createBench bool) error

	// Destroy the instances
	Destroy() error

	// Get the IP address of the LLM instance
	GetLLMInstanceIP() (string, error)

	// Get the IP address of the CPU instance
	GetBenchInstanceIP() (string, error)
}
