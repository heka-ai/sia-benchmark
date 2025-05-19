# Commands Reference

This document provides a detailed description of all available commands in the Benchmark CLI.

## Global Flags

The following flags can be used with any command:

| Flag       | Short | Description                    | Default      |
| ---------- | ----- | ------------------------------ | ------------ |
| `--config` | `-c`  | Path to the configuration file | `bench.toml` |

## Available Commands

### Root Command

```
bench
```

The root command displays help and usage information.

### Validate Configuration

```
bench validate
```

Validates the configuration file to ensure all required parameters are set correctly and the configuration follows the expected schema.

**Usage examples:**

```bash
# Validate with default config file (bench.toml)
bench validate

# Validate with a specific config file
bench validate --config my-config.toml
```

### Validate Cloud Credentials

```
bench creds
```

Validates the cloud provider credentials to ensure they are correctly set up and the CLI can connect to the cloud provider API.

**Usage examples:**

```bash
# Validate credentials with default config file
bench creds

# Validate credentials with a specific config file
bench creds --config my-config.toml
```

### Instance Management

The `instance` command group manages the cloud instances.

#### Create Instances

```
bench create
```

Creates the necessary cloud instances for running the benchmark.

**Usage examples:**

```bash
# Create instances with default config
bench create

# Create instances with a specific config file
bench create --config my-config.toml
```

### Connection Validation

```
bench connection
```

Validates the connection to the created cloud instances to ensure they are reachable and operational.

**Usage examples:**

```bash
# Validate connection with default config
bench connection

# Validate connection with a specific config file
bench connection --config my-config.toml
```

### Model Deployment

The `model` command group manages model deployment.

#### Deploy Model

```
bench deploy
```

Deploys the specified model to the cloud instance.

**Usage examples:**

```bash
# Deploy model with default config
bench deploy

# Deploy model with a specific config file
bench deploy --config my-config.toml
```

### Run Benchmark

```
bench run
```

Runs the benchmark on the deployed model.

**Usage examples:**

```bash
# Run benchmark with default config
bench run

# Run benchmark with a specific config file
bench run --config my-config.toml
```

### View Results

```
bench results
```

Displays the results of the benchmark.

**Usage examples:**

```bash
# View results with default config
bench results

# View results with a specific config file
bench results --config my-config.toml
```

### Destroy Resources

```
bench destroy
```

Destroys all cloud resources created for the benchmark.

**Usage examples:**

```bash
# Destroy resources with default config
bench destroy

# Destroy resources with a specific config file
bench destroy --config my-config.toml
```

## Command Execution Flow

The typical flow of commands for a complete benchmark session:

1. `bench validate` - Validate the configuration
2. `bench creds` - Validate cloud credentials
3. `bench create` - Create cloud instances
4. `bench connection` - Validate connection to instances
5. `bench deploy` - Deploy the model
6. `bench run` - Run the benchmark
7. `bench results` - View benchmark results
8. `bench destroy` - Clean up resources

This flow ensures that all prerequisites are met before running the benchmark and that resources are properly cleaned up afterward.
