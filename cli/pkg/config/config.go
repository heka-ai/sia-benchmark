package config

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Config struct {
	BenchID         string           `mapstructure:"bench_id" validate:"required"`
	Provider        string           `mapstructure:"provider" validate:"required,oneof=aws gcp scaleway"`
	InferenceEngine string           `mapstructure:"inference_engine" validate:"required,oneof=vllm"`
	AWSConfig       *AWSConfig       `mapstructure:"aws" validate:"required_if=Provider aws"`
	VLLMConfig      *VLLMConfig      `mapstructure:"vllm" validate:"required_if=InferenceEngine vllm"`
	InstanceConfig  *InstanceConfig  `mapstructure:"instance"`
	BenchmarkConfig *BenchmarkConfig `mapstructure:"benchmark" validate:"required"`
	APIKey          string           `mapstructure:"api_key" validate:"required"`
}

type BenchmarkConfig struct {
	Token       string `mapstructure:"token" json:"token" validate:"required"`
	DatasetName string `mapstructure:"dataset_name" json:"dataset-name" validate:"required"`
	DatasetPath string `mapstructure:"dataset_path" json:"dataset-path" validate:"required"`
	HFRevision  string `mapstructure:"hf_revision" json:"hf-revision" validate:"required"`
	HFSplit     string `mapstructure:"hf_split" json:"hf-split" validate:"required"`
	NumPrompts  int    `mapstructure:"num_prompts" json:"num-prompts" validate:"required"`
	Seed        int    `mapstructure:"seed" json:"seed" validate:"required"`
}

type VLLMConfig struct {
	Model                                      string   `mapstructure:"model"  validate:"required"`
	Task                                       *string  `mapstructure:"task" json:"task" validate:"omitempty,oneof=auto generate embedding embed classify score reward"`
	Tokenizer                                  *string  `mapstructure:"tokenizer" json:"tokenizer" validate:"omitempty"`
	SkipTokenizerInit                          *bool    `mapstructure:"skip-tokenizer-init" json:"skip-tokenizer-init"`
	Revision                                   *string  `mapstructure:"revision" json:"revision" validate:"omitempty"`
	CodeRevision                               *string  `mapstructure:"code-revision" json:"code-revision" validate:"omitempty"`
	TokenizerRevision                          *string  `mapstructure:"tokenizer-revision" json:"tokenizer-revision" validate:"omitempty"`
	TokenizerMode                              *string  `mapstructure:"tokenizer-mode" json:"tokenizer-mode" validate:"omitempty,oneof=auto slow mistral"`
	TrustRemoteCode                            *bool    `mapstructure:"trust-remote-code" json:"trust-remote-code"`
	AllowedLocalMediaPath                      *string  `mapstructure:"allowed-local-media-path" json:"allowed-local-media-path" validate:"omitempty"`
	DownloadDir                                *string  `mapstructure:"download-dir" json:"download-dir" validate:"omitempty"`
	LoadFormat                                 *string  `mapstructure:"load-format" json:"load-format" validate:"omitempty,oneof=auto pt safetensors npcache dummy tensorizer sharded-state gguf bitsandbytes mistral runai-streamer"`
	ConfigFormat                               *string  `mapstructure:"config-format" json:"config-format" validate:"omitempty,oneof=auto hf mistral"`
	Dtype                                      *string  `mapstructure:"dtype" json:"dtype" validate:"omitempty,oneof=auto half float16 bfloat16 float float32"`
	KVCacheDtype                               *string  `mapstructure:"kv-cache-dtype" json:"kv-cache-dtype" validate:"omitempty,oneof=auto fp8 fp8-e5m2 fp8-e4m3"`
	MaxModelLen                                *int     `mapstructure:"max-model-len" json:"max-model-len" validate:"omitempty"`
	GuidedDecodingBackend                      *string  `mapstructure:"guided-decoding-backend" json:"guided-decoding-backend" validate:"omitempty,oneof=outlines lm-format-enforcer xgrammar"`
	LogitsProcessorPattern                     *string  `mapstructure:"logits-processor-pattern" json:"logits-processor-pattern" validate:"omitempty"`
	ModelImpl                                  *string  `mapstructure:"model-impl" json:"model-impl" validate:"omitempty,oneof=auto vllm transformers"`
	DistributedExecutorBackend                 *string  `mapstructure:"distributed-executor-backend" json:"distributed-executor-backend" validate:"omitempty,oneof=ray mp uni external-launcher"`
	PipelineParallelSize                       *int     `mapstructure:"pipeline-parallel-size" json:"pipeline-parallel-size" validate:"omitempty"`
	TensorParallelSize                         *int     `mapstructure:"tensor-parallel-size" json:"tensor-parallel-size" validate:"omitempty"`
	MaxParallelLoadingWorkers                  *int     `mapstructure:"max-parallel-loading-workers" json:"max-parallel-loading-workers" validate:"omitempty"`
	RayWorkersUseNsight                        *bool    `mapstructure:"ray-workers-use-nsight" json:"ray-workers-use-nsight"`
	BlockSize                                  *int     `mapstructure:"block-size" json:"block-size" validate:"omitempty,oneof=8 16 32 64 128"`
	EnablePrefixCaching                        *bool    `mapstructure:"enable-prefix-caching" json:"enable-prefix-caching"`
	DisableSlidingWindow                       *bool    `mapstructure:"disable-sliding-window" json:"disable-sliding-window"`
	UseV2BlockManager                          *bool    `mapstructure:"use-v2-block-manager" json:"use-v2-block-manager"`
	NumLookaheadSlots                          *int     `mapstructure:"num-lookahead-slots" json:"num-lookahead-slots" validate:"omitempty"`
	Seed                                       *int     `mapstructure:"seed" json:"seed" validate:"omitempty"`
	SwapSpace                                  *int     `mapstructure:"swap-space" json:"swap-space" validate:"omitempty"`
	CPUOffloadGB                               *int     `mapstructure:"cpu-offload-gb" json:"cpu-offload-gb" validate:"omitempty"`
	GPUMemoryUtilization                       *int     `mapstructure:"gpu-memory-utilization" json:"gpu-memory-utilization" validate:"omitempty"`
	NumGPUBlocksOverride                       *int     `mapstructure:"num-gpu-blocks-override" json:"num-gpu-blocks-override" validate:"omitempty"`
	MaxNumBatchedTokens                        *int     `mapstructure:"max-num-batched-tokens" json:"max-num-batched-tokens" validate:"omitempty"`
	MaxNumSeqs                                 *int     `mapstructure:"max-num-seqs" json:"max-num-seqs" validate:"omitempty"`
	MaxLogprobs                                *int     `mapstructure:"max-logprobs" json:"max-logprobs" validate:"omitempty"`
	DisableLogStats                            *bool    `mapstructure:"disable-log-stats" json:"disable-log-stats"`
	Quantization                               *string  `mapstructure:"quantization" json:"quantization" validate:"omitempty,oneof=aqlm awq deepspeedfp tpu-int8 fp8 fbgemm-fp8 modelopt marlin gguf gptq-marlin-24 gptq-marlin awq-marlin gptq compressed-tensors bitsandbytes qqq hqq experts-int8 neuron-quant ipex quark moe-wna16 None"`
	RopeScaling                                *string  `mapstructure:"rope-scaling" json:"rope-scaling" validate:"omitempty"`
	RopeTheta                                  *string  `mapstructure:"rope-theta" json:"rope-theta" validate:"omitempty"`
	HFOverrides                                *string  `mapstructure:"hf-overrides" json:"hf-overrides" validate:"omitempty"`
	EnforceEager                               *bool    `mapstructure:"enforce-eager" json:"enforce-eager"`
	MaxSeqLenToCapture                         *int     `mapstructure:"max-seq-len-to-capture" json:"max-seq-len-to-capture" validate:"omitempty"`
	DisableCustomAllReduce                     *bool    `mapstructure:"disable-custom-all-reduce" json:"disable-custom-all-reduce"`
	TokenizerPoolSize                          *int     `mapstructure:"tokenizer-pool-size" json:"tokenizer-pool-size" validate:"omitempty"`
	TokenizerPoolType                          *string  `mapstructure:"tokenizer-pool-type" json:"tokenizer-pool-type" validate:"omitempty"`
	TokenizerPoolExtraConfig                   *string  `mapstructure:"tokenizer-pool-extra-config" json:"tokenizer-pool-extra-config" validate:"omitempty"`
	LimitMMPerPrompt                           *int     `mapstructure:"limit-mm-per-prompt" json:"limit-mm-per-prompt" validate:"omitempty"`
	MMProcessorKwargs                          *string  `mapstructure:"mm-processor-kwargs" json:"mm-processor-kwargs" validate:"omitempty"`
	DisableMMPreprocessorCache                 *bool    `mapstructure:"disable-mm-preprocessor-cache" json:"disable-mm-preprocessor-cache"`
	EnableLora                                 *bool    `mapstructure:"enable-lora" json:"enable-lora"`
	EnableLoraBias                             *bool    `mapstructure:"enable-lora-bias" json:"enable-lora-bias"`
	MaxLoras                                   *int     `mapstructure:"max-loras" json:"max-loras" validate:"omitempty"`
	MaxLoraRank                                *int     `mapstructure:"max-lora-rank" json:"max-lora-rank" validate:"omitempty"`
	LoraExtraVocabSize                         *int     `mapstructure:"lora-extra-vocab-size" json:"lora-extra-vocab-size" validate:"omitempty"`
	LoraDtype                                  *string  `mapstructure:"lora-dtype" json:"lora-dtype" validate:"omitempty,oneof=auto float16 bfloat16"`
	LongLoraScalingFactors                     *string  `mapstructure:"long-lora-scaling-factors" json:"long-lora-scaling-factors" validate:"omitempty"`
	MaxCPULoras                                *int     `mapstructure:"max-cpu-loras" json:"max-cpu-loras" validate:"omitempty"`
	FullyShardedLoras                          *bool    `mapstructure:"fully-sharded-loras" json:"fully-sharded-loras"`
	EnablePromptAdapter                        *bool    `mapstructure:"enable-prompt-adapter" json:"enable-prompt-adapter"`
	MaxPromptAdapters                          *int     `mapstructure:"max-prompt-adapters" json:"max-prompt-adapters" validate:"omitempty"`
	MaxPromptAdapterToken                      *int     `mapstructure:"max-prompt-adapter-token" json:"max-prompt-adapter-token" validate:"omitempty"`
	Device                                     *string  `mapstructure:"device" json:"device" validate:"omitempty,oneof=auto cuda neuron cpu openvino tpu xpu hpu"`
	NumSchedulerSteps                          *int     `mapstructure:"num-scheduler-steps" json:"num-scheduler-steps" validate:"omitempty"`
	MultiStepStreamOutputs                     *bool    `mapstructure:"multi-step-stream-outputs" json:"multi-step-stream-outputs"`
	SchedulerDelayFactor                       *int     `mapstructure:"scheduler-delay-factor" json:"scheduler-delay-factor" validate:"omitempty"`
	EnableChunkedPrefill                       *bool    `mapstructure:"enable-chunked-prefill" json:"enable-chunked-prefill"`
	SpeculativeModel                           *string  `mapstructure:"speculative-model" json:"speculative-model" validate:"omitempty"`
	SpeculativeModelQuantization               *string  `mapstructure:"speculative-model-quantization" json:"speculative-model-quantization" validate:"omitempty,oneof=aqlm awq deepspeedfp tpu-int8 fp8 fbgemm-fp8 modelopt marlin gguf gptq-marlin-24 gptq-marlin awq-marlin gptq compressed-tensors bitsandbytes qqq hqq experts-int8 neuron-quant ipex quark moe-wna16 None"`
	NumSpeculativeTokens                       *int     `mapstructure:"num-speculative-tokens" json:"num-speculative-tokens" validate:"omitempty"`
	SpeculativeDisableMQAScorer                *bool    `mapstructure:"speculative-disable-mqa-scorer" json:"speculative-disable-mqa-scorer"`
	SpeculativeDraftTensorParallelSize         *int     `mapstructure:"speculative-draft-tensor-parallel-size" json:"speculative-draft-tensor-parallel-size" validate:"omitempty"`
	SpeculativeMaxModelLen                     *int     `mapstructure:"speculative-max-model-len" json:"speculative-max-model-len" validate:"omitempty"`
	SpeculativeDisableByBatchSize              *bool    `mapstructure:"speculative-disable-by-batch-size" json:"speculative-disable-by-batch-size"`
	NgramPromptLookupMax                       *int     `mapstructure:"ngram-prompt-lookup-max" json:"ngram-prompt-lookup-max" validate:"omitempty"`
	NgramPromptLookupMin                       *int     `mapstructure:"ngram-prompt-lookup-min" json:"ngram-prompt-lookup-min" validate:"omitempty"`
	SpecDecodingAcceptanceMethod               *string  `mapstructure:"spec-decoding-acceptance-method" json:"spec-decoding-acceptance-method" validate:"omitempty,oneof=rejection-sampler typical-acceptance-sampler"`
	TypicalAcceptanceSamplerPosteriorThreshold *float64 `mapstructure:"typical-acceptance-sampler-posterior-threshold" json:"typical-acceptance-sampler-posterior-threshold" validate:"omitempty"`
	TypicalAcceptanceSamplerPosteriorAlpha     *float64 `mapstructure:"typical-acceptance-sampler-posterior-alpha" json:"typical-acceptance-sampler-posterior-alpha" validate:"omitempty"`
	DisableLogprobsDuringSpecDecoding          *bool    `mapstructure:"disable-logprobs-during-spec-decoding" json:"disable-logprobs-during-spec-decoding"`
	ModelLoaderExtraConfig                     *string  `mapstructure:"model-loader-extra-config" json:"model-loader-extra-config" validate:"omitempty"`
	IgnorePatterns                             *string  `mapstructure:"ignore-patterns" json:"ignore-patterns" validate:"omitempty"`
	PreemptionMode                             *string  `mapstructure:"preemption-mode" json:"preemption-mode" validate:"omitempty"`
	QLoraAdapterNameOrPath                     *string  `mapstructure:"qlora-adapter-name-or-path" json:"qlora-adapter-name-or-path" validate:"omitempty"`
	OtlpTracesEndpoint                         *string  `mapstructure:"otlp-traces-endpoint" json:"otlp-traces-endpoint" validate:"omitempty"`
	CollectDetailedTraces                      *bool    `mapstructure:"collect-detailed-traces" json:"collect-detailed-traces"`
	DisableAsyncOutputProc                     *bool    `mapstructure:"disable-async-output-proc" json:"disable-async-output-proc"`
	SchedulingPolicy                           *string  `mapstructure:"scheduling-policy" json:"scheduling-policy" validate:"omitempty,oneof=fcfs priority"`
	OverrideNeuronConfig                       *string  `mapstructure:"override-neuron-config" json:"override-neuron-config" validate:"omitempty"`
	OverridePoolerConfig                       *string  `mapstructure:"override-pooler-config" json:"override-pooler-config" validate:"omitempty"`
	CompilationConfig                          *string  `mapstructure:"compilation-config" json:"compilation-config" validate:"omitempty"`
	KVTransferConfig                           *string  `mapstructure:"kv-transfer-config" json:"kv-transfer-config" validate:"omitempty"`
	WorkerCls                                  *string  `mapstructure:"worker-cls" json:"worker-cls" validate:"omitempty"`
	GenerationConfig                           *string  `mapstructure:"generation-config" json:"generation-config" validate:"omitempty"`
	OverrideGenerationConfig                   *string  `mapstructure:"override-generation-config" json:"override-generation-config" validate:"omitempty"`
	EnableSleepMode                            *bool    `mapstructure:"enable-sleep-mode" json:"enable-sleep-mode"`
	CalculateKVScales                          *bool    `mapstructure:"calculate-kv-scales" json:"calculate-kv-scales"`
}

type AWSConfig struct {
	Region          string `mapstructure:"region" validate:"required"`
	CPUInstanceType string `mapstructure:"cpu_instance_type" validate:"required"`
	GPUInstanceType string `mapstructure:"gpu_instance_type" validate:"required"`

	AWSAccessKey string `mapstructure:"access_key" validate:"required_if=ProfileName false"`
	AWSSecretKey string `mapstructure:"secret_key" validate:"required_if=ProfileName false"`

	ProfileName string `mapstructure:"profile_name" validate:"required_if=AWSAccessKey false"`

	GPU_AMI string `mapstructure:"gpu_ami" validate:"required"`
	CPU_AMI string `mapstructure:"cpu_ami" validate:"required"`
}

type GCPConfig struct {
	Region          string `mapstructure:"region" validate:"required"`
	CPUInstanceType string `mapstructure:"cpu_instance_type" validate:"required"`
	GPUInstanceType string `mapstructure:"gpu_instance_type" validate:"required"`

	GCPAccessKey string `mapstructure:"access_key" validate:"required"`
	GCPSecretKey string `mapstructure:"secret_key" validate:"required"`
}

type ScalewayConfig struct {
	Region          string `mapstructure:"region" validate:"required"`
	CPUInstanceType string `mapstructure:"cpu_instance_type" validate:"required"`
	GPUInstanceType string `mapstructure:"gpu_instance_type" validate:"required"`

	ScalewayAccessKey string `mapstructure:"access_key" validate:"required"`
	ScalewaySecretKey string `mapstructure:"secret_key" validate:"required"`
}

type InstanceConfig struct {
	Test        *string `mapstructure:"test"`
	HealthCheck *string `mapstructure:"health_check"`
	SecondTest  *int    `mapstructure:"second_test"`
}

func GenerateVLLMCommand(vllmConfig *VLLMConfig) ([]string, error) {
	localArgs := []string{"serve", vllmConfig.Model}

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(vllmConfig)
	json.Unmarshal(inrec, &inInterface)

	for k, v := range inInterface {
		if k == "Model" {
			continue
		}

		if v == nil || v == "" {
			continue
		}

		s, ok := v.(string)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), s)
			continue
		}

		number, ok := v.(float64)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatFloat(number, 'f', 0, 64))
			continue
		}

		b, ok := v.(bool)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatBool(b))
			continue
		}

		logger.Warn().Str("key", k).Interface("value", v).Msg("Unknown type")
	}

	return localArgs, nil
}

func GenerateBenchmarkCommand(conf *Config, ip string) ([]string, error) {
	localArgs := []string{"/home/ubuntu/ec2/cpu/benchmark.py", fmt.Sprintf("--base-url http://%s:8000", ip)}

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(conf.BenchmarkConfig)
	json.Unmarshal(inrec, &inInterface)

	for k, v := range inInterface {
		if v == nil || v == "" {
			continue
		}

		if k == "token" {
			continue
		}

		s, ok := v.(string)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), s)
			continue
		}

		number, ok := v.(float64)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatFloat(number, 'f', 0, 64))
			continue
		}

		b, ok := v.(bool)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatBool(b))
			continue
		}

		logger.Warn().Str("key", k).Interface("value", v).Msg("Unknown type")
	}

	localArgs = append(localArgs, "--model", conf.VLLMConfig.Model)

	return localArgs, nil
}
