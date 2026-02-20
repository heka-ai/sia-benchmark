package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GeneralConfig   *GeneralConfig   `mapstructure:"general" validate:"required"`
	AWSConfig       *AWSConfig       `mapstructure:"aws"`
	LocalConfig     *LocalConfig     `mapstructure:"local"`
	VLLMConfig      *VLLMConfig      `mapstructure:"vllm"`
	BenchmarkConfig *BenchmarkConfig `mapstructure:"benchmark" validate:"required"`
}

type GeneralConfig struct {
	BenchmarkID     string `mapstructure:"benchmark_id" validate:"required"`
	BenchmarkName   string `mapstructure:"benchmark_name"`
	Provider        string `mapstructure:"provider" validate:"required,oneof=aws gcp scaleway local"`
	InferenceEngine string `mapstructure:"inference_engine" validate:"required,oneof=vllm"`
	APIKey          string `mapstructure:"api_key" validate:"required_unless=Provider local"`
}

type BenchmarkConfig struct {
	Token          string `mapstructure:"token" json:"token" validate:"required"`
	DatasetName    string `mapstructure:"dataset_name" json:"dataset-name" validate:"required,oneof=random hf heka"`
	DatasetPath    string `mapstructure:"dataset_path" json:"dataset-path" validate:"required"`
	HFRevision     string `mapstructure:"hf_revision" json:"hf-revision" validate:"required_unless=DatasetName random" default:"main"`
	HFSplit        string `mapstructure:"hf_split" json:"hf-split" validate:"required_unless=DatasetName random" default:"train"`
	NumPrompts     int    `mapstructure:"num_prompts" json:"num-prompts" validate:"required" default:"500"`
	Seed           int    `mapstructure:"seed" json:"seed" validate:"required" default:"42"`
	Backend        string `mapstructure:"backend" json:"backend" validate:"omitempty,oneof=openai" default:"openai"`
	SaveResult     bool   `mapstructure:"save_result" json:"save-result" validate:"omitempty" default:"true"`
	ResultFilename string `mapstructure:"result_filename" json:"result-filename" validate:"omitempty" default:"metrics.json"`
	EnginePort     int    `mapstructure:"engine_port" json:"engine-port" validate:"omitempty" default:"8000"`
	ScriptPath     string `mapstructure:"script_path" json:"script-path" validate:"omitempty"`
}

type VLLMConfig struct {
	Model                                      string   `mapstructure:"model"  validate:"required"`
	Task                                       *string  `mapstructure:"task" json:"task" validate:"omitempty,oneof=None auto generate embedding embed classify score reward"`
	Tokenizer                                  *string  `mapstructure:"tokenizer" json:"tokenizer" validate:"omitempty"`
	SkipTokenizerInit                          *bool    `mapstructure:"skip-tokenizer-init" json:"skip-tokenizer-init"`
	Revision                                   *string  `mapstructure:"revision" json:"revision" validate:"omitempty"`
	CodeRevision                               *string  `mapstructure:"code-revision" json:"code-revision" validate:"omitempty"`
	TokenizerRevision                          *string  `mapstructure:"tokenizer-revision" json:"tokenizer-revision" validate:"omitempty"`
	TokenizerMode                              *string  `mapstructure:"tokenizer-mode" json:"tokenizer-mode" validate:"omitempty,oneof=auto slow mistral"`
	TrustRemoteCode                            *bool    `mapstructure:"trust-remote-code" json:"trust-remote-code"`
	AllowedLocalMediaPath                      *string  `mapstructure:"allowed-local-media-path" json:"allowed-local-media-path" validate:"omitempty"`
	DownloadDir                                *string  `mapstructure:"download-dir" json:"download-dir" validate:"omitempty"`
	LoadFormat                                 *string  `mapstructure:"load-format" json:"load-format" validate:"omitempty,oneof=None auto pt safetensors npcache dummy tensorizer sharded-state gguf bitsandbytes mistral runai-streamer"`
	ConfigFormat                               *string  `mapstructure:"config-format" json:"config-format" validate:"omitempty,oneof=None auto hf mistral"`
	Dtype                                      *string  `mapstructure:"dtype" json:"dtype" validate:"omitempty,oneof=auto half float16 bfloat16 float float32"`
	KVCacheDtype                               *string  `mapstructure:"kv-cache-dtype" json:"kv-cache-dtype" validate:"omitempty,oneof=auto fp8 fp8-e5m2 fp8-e4m3"`
	MaxModelLen                                *int     `mapstructure:"max-model-len" json:"max-model-len" validate:"omitempty"`
	GuidedDecodingBackend                      *string  `mapstructure:"guided-decoding-backend" json:"guided-decoding-backend" validate:"omitempty,oneof=None outlines lm-format-enforcer xgrammar"`
	LogitsProcessorPattern                     *string  `mapstructure:"logits-processor-pattern" json:"logits-processor-pattern" validate:"omitempty"`
	ModelImpl                                  *string  `mapstructure:"model-impl" json:"model-impl" validate:"omitempty,oneof=None auto vllm transformers"`
	DistributedExecutorBackend                 *string  `mapstructure:"distributed-executor-backend" json:"distributed-executor-backend" validate:"omitempty,oneof=None ray mp uni external-launcher"`
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
	GPUMemoryUtilization                       *float64 `mapstructure:"gpu-memory-utilization" json:"gpu-memory-utilization" validate:"omitempty"`
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
	LoraDtype                                  *string  `mapstructure:"lora-dtype" json:"lora-dtype" validate:"omitempty,oneof=None auto float16 bfloat16"`
	LongLoraScalingFactors                     *string  `mapstructure:"long-lora-scaling-factors" json:"long-lora-scaling-factors" validate:"omitempty"`
	MaxCPULoras                                *int     `mapstructure:"max-cpu-loras" json:"max-cpu-loras" validate:"omitempty"`
	FullyShardedLoras                          *bool    `mapstructure:"fully-sharded-loras" json:"fully-sharded-loras"`
	EnablePromptAdapter                        *bool    `mapstructure:"enable-prompt-adapter" json:"enable-prompt-adapter"`
	MaxPromptAdapters                          *int     `mapstructure:"max-prompt-adapters" json:"max-prompt-adapters" validate:"omitempty"`
	MaxPromptAdapterToken                      *int     `mapstructure:"max-prompt-adapter-token" json:"max-prompt-adapter-token" validate:"omitempty"`
	Device                                     *string  `mapstructure:"device" json:"device" validate:"omitempty,oneof=None auto cuda neuron cpu openvino tpu xpu hpu"`
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
	SpecDecodingAcceptanceMethod               *string  `mapstructure:"spec-decoding-acceptance-method" json:"spec-decoding-acceptance-method" validate:"omitempty,oneof=None rejection-sampler typical-acceptance-sampler"`
	TypicalAcceptanceSamplerPosteriorThreshold *float64 `mapstructure:"typical-acceptance-sampler-posterior-threshold" json:"typical-acceptance-sampler-posterior-threshold" validate:"omitempty"`
	TypicalAcceptanceSamplerPosteriorAlpha     *float64 `mapstructure:"typical-acceptance-sampler-posterior-alpha" json:"typical-acceptance-sampler-posterior-alpha" validate:"omitempty"`
	DisableLogprobsDuringSpecDecoding          *bool    `mapstructure:"disable-logprobs-during-spec-decoding" json:"disable-logprobs-during-spec-decoding"`
	ModelLoaderExtraConfig                     *string  `mapstructure:"model-loader-extra-config" json:"model-loader-extra-config" validate:"omitempty"`
	IgnorePatterns                             []string `mapstructure:"ignore-patterns" json:"ignore-patterns" validate:"omitempty"`
	PreemptionMode                             *string  `mapstructure:"preemption-mode" json:"preemption-mode" validate:"omitempty"`
	QLoraAdapterNameOrPath                     *string  `mapstructure:"qlora-adapter-name-or-path" json:"qlora-adapter-name-or-path" validate:"omitempty"`
	OtlpTracesEndpoint                         *string  `mapstructure:"otlp-traces-endpoint" json:"otlp-traces-endpoint" validate:"omitempty"`
	CollectDetailedTraces                      *bool    `mapstructure:"collect-detailed-traces" json:"collect-detailed-traces"`
	DisableAsyncOutputProc                     *bool    `mapstructure:"disable-async-output-proc" json:"disable-async-output-proc"`
	SchedulingPolicy                           *string  `mapstructure:"scheduling-policy" json:"scheduling-policy" validate:"omitempty,oneof=None fcfs priority"`
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

	AMI string `mapstructure:"ami" validate:"required"`
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

type LocalConfig struct {
	AcceleratorType string `mapstructure:"accelerator_type" validate:"required,oneof=jetson nvidia"`
	GPUInstanceType string `mapstructure:"gpu_instance_type" validate:"required"`
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
			// Skip explicitly disabled or unset string options
			if s == "None" {
				continue
			}
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), s)
			continue
		}

		if number, ok := v.(float64); ok {
			// preserve decimals for float flags
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatFloat(number, 'f', -1, 64))
			continue
		}

		b, ok := v.(bool)
		if ok {
			// Only include boolean flags when explicitly true
			if b {
				localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatBool(b))
			}
			continue
		}

		logger.Warn().Str("key", k).Interface("value", v).Msg("Unknown type")
	}

	// Use benchmark engine_port when available; otherwise default to 8000.
	enginePort := 8000
	cfg := GetConfig()
	if cfg.BenchmarkConfig != nil && cfg.BenchmarkConfig.EnginePort > 0 {
		enginePort = cfg.BenchmarkConfig.EnginePort
	}
	localArgs = append(localArgs, "--port", strconv.Itoa(enginePort))

	return localArgs, nil
}

func GenerateBenchmarkCommand(conf *Config, ip string, port int) ([]string, error) {
	script := resolveBenchmarkScriptPath(conf)
	if strings.TrimSpace(script) == "" {
		return nil, fmt.Errorf("benchmark.py not found; set benchmark.script_path or BENCHMARK_SCRIPT")
	}

	// Compute effective backend: fallback to top-level InferenceEngine when unset
	effectiveBackend := strings.TrimSpace(conf.BenchmarkConfig.Backend)
	if effectiveBackend == "" {
		effectiveBackend = strings.TrimSpace(conf.GeneralConfig.InferenceEngine)
	}

	localArgs := []string{script, "--backend", effectiveBackend, "--base-url", fmt.Sprintf("http://%s:%d", ip, port)}

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(conf.BenchmarkConfig)
	json.Unmarshal(inrec, &inInterface)

	for k, v := range inInterface {
		if v == nil || v == "" {
			continue
		}

		// Exclude fields not meant to be passed to benchmark.py directly
		if k == "token" || k == "backend" || k == "script_path" || k == "engine-port" {
			continue
		}

		s, ok := v.(string)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), s)
			continue
		}

		if number, ok := v.(float64); ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k), strconv.FormatFloat(number, 'f', -1, 64))
			continue
		}

		_, ok = v.(bool)
		if ok {
			localArgs = append(localArgs, fmt.Sprintf("--%s", k))
			continue
		}

		logger.Warn().Str("key", k).Interface("value", v).Msg("Unknown type")
	}

	localArgs = append(localArgs, "--model", conf.VLLMConfig.Model)

	// Forward the CLI --config path as the deployment config file
	if cfgPath := strings.TrimSpace(viper.GetString("config")); cfgPath != "" {
		localArgs = append(localArgs, "--deployment-config-file", cfgPath)
	}

	return localArgs, nil
}

// resolveBenchmarkScriptPath finds the benchmark.py path in a portable way:
// 1) BENCHMARK_SCRIPT env var
// 2) benchmark.script_path in TOML
// 3) common repo-relative locations
func resolveBenchmarkScriptPath(conf *Config) string {
	// 1) env var
	if s := strings.TrimSpace(os.Getenv("BENCHMARK_SCRIPT")); s != "" {
		return s
	}
	// 2) from config
	if conf != nil && conf.BenchmarkConfig != nil {
		if s := strings.TrimSpace(conf.BenchmarkConfig.ScriptPath); s != "" {
			return s
		}
	}
	// 3) candidates relative to current working dir
	candidates := []string{
		"benchmark.py",
		"./benchmark.py",
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if st, err := os.Stat(abs); err == nil && !st.IsDir() {
			return abs
		}
	}
	return ""
}
