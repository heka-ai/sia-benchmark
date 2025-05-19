# Configuration Reference

The Benchmark CLI uses a TOML configuration file to define benchmark settings. By default, it looks for a file named `bench.toml` in the current directory, but you can specify a different path using the `--config` flag.

## Top-Level Configuration

| Parameter          | Type   | Description                                         | Required |
| ------------------ | ------ | --------------------------------------------------- | -------- |
| `bench_id`         | String | Unique identifier for the benchmark run             | Yes      |
| `provider`         | String | Cloud provider to use (`aws`, `gcp`, or `scaleway`) | Yes      |
| `inference_engine` | String | Inference engine to use (`vllm`)                    | Yes      |

## Cloud Provider Configuration

### AWS Configuration

Define AWS-specific settings in the `[aws]` section:

| Parameter                | Type   | Description                                | Required |
| ------------------------ | ------ | ------------------------------------------ | -------- |
| `region`                 | String | AWS region where resources will be created | Yes      |
| `gpu_ami`                | String | AMI ID for GPU instances                   | Yes      |
| `cpu_ami`                | String | AMI ID for CPU instances                   | Yes      |
| `instance_type`          | String | Instance type for the model server         | Yes      |
| `launcher_instance_type` | String | Instance type for the benchmark runner     | Yes      |
| `profile_name`           | String | AWS profile name from your AWS credentials | Yes\*    |
| `access_key`             | String | AWS access key ID                          | Yes\*    |
| `secret_key`             | String | AWS secret access key                      | Yes\*    |

\*Note: Either `profile_name` OR both `access_key` and `secret_key` must be provided.

Example:

```toml
[aws]
region = "us-east-1"
gpu_ami = "ami-072c3e2520d9af5fa"
cpu_ami = "ami-04f3f32777c02a5b3"
instance_type = "g4dn.xlarge"
launcher_instance_type = "t3.micro"
profile_name = "my-aws-profile"
```

### GCP Configuration

Define GCP-specific settings in the `[gcp]` section:

| Parameter                | Type   | Description                                | Required |
| ------------------------ | ------ | ------------------------------------------ | -------- |
| `region`                 | String | GCP region where resources will be created | Yes      |
| `instance_type`          | String | Instance type for the model server         | Yes      |
| `launcher_instance_type` | String | Instance type for the benchmark runner     | Yes      |
| `access_key`             | String | GCP access key                             | Yes      |
| `secret_key`             | String | GCP secret key                             | Yes      |

### Scaleway Configuration

Define Scaleway-specific settings in the `[scaleway]` section:

| Parameter                | Type   | Description                                     | Required |
| ------------------------ | ------ | ----------------------------------------------- | -------- |
| `region`                 | String | Scaleway region where resources will be created | Yes      |
| `instance_type`          | String | Instance type for the model server              | Yes      |
| `launcher_instance_type` | String | Instance type for the benchmark runner          | Yes      |
| `access_key`             | String | Scaleway access key                             | Yes      |
| `secret_key`             | String | Scaleway secret key                             | Yes      |

## Inference Engine Configuration

### vLLM Configuration

Define vLLM-specific settings in the `[vllm]` section:

| Parameter           | Type    | Description                                                                         | Required |
| ------------------- | ------- | ----------------------------------------------------------------------------------- | -------- |
| `model`             | String  | Model name or path                                                                  | Yes      |
| `task`              | String  | Task type (`auto`, `generate`, `embedding`, `embed`, `classify`, `score`, `reward`) | No       |
| `tokenizer`         | String  | Custom tokenizer to use                                                             | No       |
| `trust_remote_code` | Boolean | Whether to trust remote code                                                        | No       |
| `quantization`      | String  | Quantization method                                                                 | No       |
| `dtype`             | String  | Data type for model weights                                                         | No       |
| `max_model_len`     | Integer | Maximum sequence length                                                             | No       |

_Note: This is a subset of the available vLLM parameters. For a complete list, refer to the internal/config/config.go file._

Example:

```toml
[vllm]
model = "mistralai/Mistral-7B-v0.1"
task = "generate"
trust_remote_code = true
dtype = "bfloat16"
max_model_len = 4096
```

## Instance Configuration

Define instance-related settings in the `[instance]` section:

| Parameter      | Type   | Description                     | Required |
| -------------- | ------ | ------------------------------- | -------- |
| `health_check` | String | HTTP endpoint for health checks | Yes      |

Example:

```toml
[instance]
health_check = "/health"
```

## Full Configuration Example

```toml
bench_id = "my-benchmark"
provider = "aws"
inference_engine = "vllm"

[aws]
region = "us-east-1"
gpu_ami = "ami-072c3e2520d9af5fa"
cpu_ami = "ami-04f3f32777c02a5b3"
instance_type = "g4dn.xlarge"
launcher_instance_type = "t3.micro"
profile_name = "my-aws-profile"

[vllm]
model = "mistralai/Mistral-7B-v0.1"
task = "generate"
trust_remote_code = true
dtype = "bfloat16"
max_model_len = 4096

[instance]
health_check = "/health"
```
