# Getting Started with Benchmark CLI

This guide will help you get started with the Benchmark CLI tool for LLM benchmarking.

## Installation

### Prerequisites

- Go 1.20 or later
- Access to a supported cloud provider (AWS, GCP, or Scaleway)
- Configured cloud provider credentials

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/heka-ai/benchmark-cli.git
   cd benchmark-cli
   ```

2. Build the binary:

   ```bash
   go build -o bench cmd/bench.go
   ```

3. Move the binary to your PATH (optional):
   ```bash
   sudo mv bench /usr/local/bin/
   ```

## Basic Usage

### Configuration

Before running benchmarks, you need to create a configuration file. By default, the CLI looks for a file named `bench.toml` in the current directory.

Here's a basic configuration example:

```toml
bench_id = "my-benchmark"
provider = "aws"
inference_engine = "vllm"

[aws]
region = "us-east-1"
gpu_ami = "ami-072c3e2520d9af5fa"  # AMI with GPU support
cpu_ami = "ami-04f3f32777c02a5b3"  # Standard CPU AMI
instance_type = "g4dn.xlarge"      # GPU instance for model
launcher_instance_type = "t3.micro" # Instance for sending requests

# Use either profile_name or access_key/secret_key
profile_name = "my-aws-profile"
# access_key = "YOUR_AWS_ACCESS_KEY"
# secret_key = "YOUR_AWS_SECRET_KEY"

[vllm]
model = "mistralai/Mistral-7B-v0.1"

[instance]
health_check = "/health"
```

### Validating Configuration

Validate your configuration:

```bash
bench validate
```

### Validating Cloud Credentials

Verify that your cloud credentials are correctly set up:

```bash
bench creds
```

### Creating Infrastructure

Create the necessary infrastructure:

```bash
bench instance create
```

### Checking Connection

Verify connectivity to the created instances:

```bash
bench connection
```

### Deploying a Model

Deploy the specified model:

```bash
bench model deploy
```

### Running Benchmarks

Run the benchmarks:

```bash
bench run
```

### Viewing Results

View benchmark results:

```bash
bench results
```

### Cleaning Up

When finished, clean up all created resources:

```bash
bench destroy
```

## Next Steps

- Read the [Configuration Guide](configuration.md) for detailed configuration options
- Check the [Commands Reference](commands.md) for all available commands
- Explore [Cloud Providers](cloud-providers.md) for provider-specific settings
