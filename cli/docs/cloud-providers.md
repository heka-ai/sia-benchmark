# Cloud Providers

The Benchmark CLI supports multiple cloud providers for running LLM benchmarks. This document outlines the configuration and requirements for each supported provider.

## Overview

Currently, the CLI supports the following cloud providers:

- AWS (fully implemented)
- GCP (placeholder implementation)
- Scaleway (placeholder implementation)

## AWS

AWS is the primary cloud provider supported by the Benchmark CLI. It allows you to run benchmarks on various AWS EC2 instance types.

### Requirements

- AWS account
- AWS CLI installed and configured, or AWS access and secret keys
- Sufficient AWS permissions to create and manage EC2 instances

### Configuration

Example AWS configuration:

```toml
provider = "aws"

[aws]
region = "us-east-1"
gpu_ami = "ami-072c3e2520d9af5fa"
cpu_ami = "ami-04f3f32777c02a5b3"
instance_type = "g4dn.xlarge"
launcher_instance_type = "t3.micro"

# Authentication using AWS profile
profile_name = "my-aws-profile"

# Or authentication using access keys
# access_key = "YOUR_AWS_ACCESS_KEY"
# secret_key = "YOUR_AWS_SECRET_KEY"
```

### Instance Types

For running LLM benchmarks, you typically need:

- A GPU instance for the model server (e.g., `g4dn.xlarge`, `p3.2xlarge`)
- A smaller CPU instance for the benchmark runner (e.g., `t3.micro`)

### Authentication Methods

AWS supports two authentication methods:

1. **AWS Profile**: Use the `profile_name` parameter to specify an AWS profile name from your AWS credentials file.

   ```toml
   profile_name = "my-aws-profile"
   ```

2. **Access Keys**: Directly specify AWS access and secret keys.
   ```toml
   access_key = "YOUR_AWS_ACCESS_KEY"
   secret_key = "YOUR_AWS_SECRET_KEY"
   ```

### AMI Selection

You need to specify AMIs for both GPU and CPU instances:

- `gpu_ami`: The AMI used for the model server (should include CUDA and other GPU dependencies)
- `cpu_ami`: The AMI used for the benchmark runner

## GCP (Google Cloud Platform)

GCP support is currently placeholder implementation in the codebase.

### Future Configuration (Placeholder)

```toml
provider = "gcp"

[gcp]
region = "us-central1"
instance_type = "n1-standard-8-nvidia-tesla-t4"
launcher_instance_type = "e2-micro"
access_key = "YOUR_GCP_ACCESS_KEY"
secret_key = "YOUR_GCP_SECRET_KEY"
```

## Scaleway

Scaleway support is currently placeholder implementation in the codebase.

### Future Configuration (Placeholder)

```toml
provider = "scaleway"

[scaleway]
region = "fr-par"
instance_type = "GPU-3070-S"
launcher_instance_type = "DEV1-S"
access_key = "YOUR_SCALEWAY_ACCESS_KEY"
secret_key = "YOUR_SCALEWAY_SECRET_KEY"
```

## Provider Architecture

The Benchmark CLI uses a pluggable provider architecture, making it easy to add new cloud providers in the future. Each provider implements the `Cloud` interface defined in the `internal/cloud` package.

### Adding New Providers

To add a new cloud provider:

1. Implement the `Cloud` interface for your provider
2. Add configuration structures in `internal/config/config.go`
3. Update the provider handling in the relevant command files

For more details, refer to the existing AWS implementation in the `internal/cloud/aws` package.
