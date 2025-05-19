# Architecture Overview

This document provides an overview of the Benchmark CLI's architecture, explaining its design patterns and code organization.

## Project Structure

The Benchmark CLI follows a standard Go project structure:

```
benchmark-cli/
├── cmd/              # Command-line entry points
│   └── bench.go      # Main entry point for the CLI
├── internal/         # Internal packages not meant for external use
│   ├── cloud/        # Cloud provider implementations
│   │   ├── aws/      # AWS implementation
│   │   ├── gcp/      # GCP implementation (placeholder)
│   │   └── scaleway/ # Scaleway implementation (placeholder)
│   ├── cmd/          # Command implementations
│   ├── config/       # Configuration handling
│   └── logs/         # Logging utilities
├── pkg/              # Packages that could be used by external applications
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```

## Design Patterns

### Command Pattern

The CLI uses the [Cobra](https://github.com/spf13/cobra) library for command handling, implementing a command pattern architecture. Each command is defined in a separate file in the `internal/cmd` directory.

### Dependency Injection

The CLI uses a form of dependency injection where configuration and services are passed to the components that need them, reducing tight coupling.

### Interface-Based Design

The cloud provider functionality is abstracted through interfaces, allowing for different implementations for each provider while maintaining a consistent API.

## Key Components

### Entry Point (cmd/bench.go)

The main entry point of the application. It initializes the root command and adds all subcommands.

```go
func main() {
    rootCmd := cmd.RootCmd()
    rootCmd.Execute()
}
```

### Root Command (internal/cmd/root.cmd.go)

Defines the root command and adds all subcommands:

```go
func RootCmd() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "bench",
        Short: "Run business oriented LLM benchmarks",
        Long:  `Sia Benchmark is the cli to run business oriented LLM benchmarks`,
    }

    rootCmd.AddCommand(ValidateCmd())
    rootCmd.AddCommand(ValidateCloudCmd())
    rootCmd.AddCommand(InstanceCmd())
    rootCmd.AddCommand(ConnectionCmd())
    rootCmd.AddCommand(ModelCmd())
    rootCmd.AddCommand(BenchCmd())
    rootCmd.AddCommand(ResultsCmd())
    rootCmd.AddCommand(DestroyCmd())

    rootCmd.Flags().StringP("config", "c", "bench.toml", "Path to the config file")

    return rootCmd
}
```

### Configuration (internal/config)

The configuration package handles reading, parsing, and validating the TOML configuration file:

- `config.go`: Defines the configuration structures
- `validate.go`: Handles configuration validation

### Cloud Providers (internal/cloud)

The `cloud` package defines the cloud provider interface and its implementations:

```go
// Cloud interface that all providers must implement
type Cloud interface {
    Init() Cloud
    ValidateCredentials() error
    CreateInstance() error
    GetBenchmarkInstances() ([]types.Instance, error)
    DeleteInstance() error
    // ... other methods
}
```

The AWS implementation in `internal/cloud/aws` provides a complete implementation of this interface.

### Command Implementations (internal/cmd)

Each command is implemented in a separate file:

- `validate.cmd.go`: Validates the configuration
- `validate-cloud.cmd.go`: Validates cloud credentials
- `instance.cmd.go`: Manages instances
- `connection.cmd.go`: Validates connections
- `model.cmd.go`: Manages model deployment
- `bench.cmd.go`: Runs benchmarks
- `results.cmd.go`: Displays results
- `destroy.go`: Cleans up resources

### Logging (internal/logs)

The logging package provides a consistent logging interface using [zerolog](https://github.com/rs/zerolog):

```go
func GetLogger(module string) zerolog.Logger {
    return GetMainLogger().With().Str("module", module).Timestamp().Logger()
}
```

## Execution Flow

1. The user executes a command: `bench <command>`
2. The main function initializes the root command and executes it
3. Cobra routes to the appropriate command handler
4. The command handler initializes the configuration
5. The command handler executes its logic, often using the cloud provider interface
6. Results are logged to the console

## Extending the CLI

### Adding a New Command

To add a new command:

1. Create a new file in `internal/cmd/`
2. Define a function that returns a `*cobra.Command`
3. Add the command to the root command in `internal/cmd/root.cmd.go`

### Adding a New Cloud Provider

To add a new cloud provider:

1. Create a new package in `internal/cloud/`
2. Implement the `Cloud` interface
3. Add provider-specific configuration in `internal/config/config.go`
4. Update provider handling in command implementations

### Adding a New Inference Engine

To add a new inference engine:

1. Create appropriate configuration structures in `internal/config/config.go`
2. Add engine-specific code in the command implementations
3. Update validation logic to handle the new engine
