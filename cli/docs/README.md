# Benchmark CLI Documentation

Welcome to the documentation for the Benchmark CLI tool. This tool is designed to automate and standardize the process of benchmarking Language Learning Models (LLMs) across different cloud providers.

## Documentation Contents

- [Getting Started](getting-started.md): Installation and basic setup
- [Configuration](configuration.md): Detailed configuration options explanation
- [Commands Reference](commands.md): List of available commands and their usage
- [Cloud Providers](cloud-providers.md): Supported cloud providers and their specific configurations
- [Inference Engines](inference-engines.md): Supported LLM inference engines
- [Architecture](architecture.md): Overview of the project's architecture and design

## Project Overview

The Benchmark CLI is a command-line tool that facilitates the deployment, configuration, and benchmarking of LLM models on cloud infrastructure. It currently supports AWS as a cloud provider and vLLM as an inference engine, with the architecture designed to be extensible for additional providers and engines in the future.

Key features:

- Automated cloud infrastructure provisioning for LLM benchmarking
- Standardized benchmark configuration through TOML files
- Support for various cloud providers (currently AWS)
- Integration with vLLM for LLM inference
- Comprehensive validation of configurations and cloud credentials
