package config

type Config struct {
	BenchID         string           `mapstructure:"bench_id" validate:"required"`
	Provider        string           `mapstructure:"provider" validate:"required,oneof=aws gcp scaleway"`
	InferenceEngine string           `mapstructure:"inference_engine" validate:"required,oneof=vllm"`
	AWSConfig       *AWSConfig       `mapstructure:"aws" validate:"required_if=Provider aws"`
	VLLMConfig      *VLLMConfig      `mapstructure:"vllm" validate:"required_if=InferenceEngine vllm"`
	InstanceConfig  *InstanceConfig  `mapstructure:"instance" validate:"required"`
	BenchmarkConfig *BenchmarkConfig `mapstructure:"benchmark" validate:"required"`
	APIKey          string           `mapstructure:"api_key" validate:"required"`
}

type BenchmarkConfig struct {
	Token       string `mapstructure:"token" validate:"required"`
	DatasetName string `mapstructure:"dataset_name" validate:"required"`
	DatasetPath string `mapstructure:"dataset_path" validate:"required"`
	HFRevision  string `mapstructure:"hf_revision" validate:"required"`
	HFSplit     string `mapstructure:"hf_split" validate:"required"`
	NumPrompts  int    `mapstructure:"num_prompts" validate:"required"`
	Seed        int    `mapstructure:"seed" validate:"required"`
}

type VLLMConfig struct {
	Model                                      string   `mapstructure:"model" validate:"required"`
	Task                                       string   `mapstructure:"task" validate:"omitempty,oneof=auto generate embedding embed classify score reward"`
	Tokenizer                                  string   `mapstructure:"tokenizer" validate:"omitempty"`
	SkipTokenizerInit                          bool     `mapstructure:"skip_tokenizer_init"`
	Revision                                   string   `mapstructure:"revision" validate:"omitempty"`
	CodeRevision                               string   `mapstructure:"code_revision" validate:"omitempty"`
	TokenizerRevision                          string   `mapstructure:"tokenizer_revision" validate:"omitempty"`
	TokenizerMode                              string   `mapstructure:"tokenizer_mode" validate:"omitempty,oneof=auto slow mistral"`
	TrustRemoteCode                            bool     `mapstructure:"trust_remote_code"`
	AllowedLocalMediaPath                      string   `mapstructure:"allowed_local_media_path" validate:"omitempty"`
	DownloadDir                                string   `mapstructure:"download_dir" validate:"omitempty"`
	LoadFormat                                 string   `mapstructure:"load_format" validate:"omitempty,oneof=auto pt safetensors npcache dummy tensorizer sharded_state gguf bitsandbytes mistral runai_streamer"`
	ConfigFormat                               string   `mapstructure:"config_format" validate:"omitempty,oneof=auto hf mistral"`
	Dtype                                      string   `mapstructure:"dtype" validate:"omitempty,oneof=auto half float16 bfloat16 float float32"`
	KVCacheDtype                               string   `mapstructure:"kv_cache_dtype" validate:"omitempty,oneof=auto fp8 fp8_e5m2 fp8_e4m3"`
	MaxModelLen                                int      `mapstructure:"max_model_len" validate:"omitempty"`
	GuidedDecodingBackend                      string   `mapstructure:"guided_decoding_backend" validate:"omitempty,oneof=outlines lm-format-enforcer xgrammar"`
	LogitsProcessorPattern                     string   `mapstructure:"logits_processor_pattern" validate:"omitempty"`
	ModelImpl                                  string   `mapstructure:"model_impl" validate:"omitempty,oneof=auto vllm transformers"`
	DistributedExecutorBackend                 string   `mapstructure:"distributed_executor_backend" validate:"omitempty,oneof=ray mp uni external_launcher"`
	PipelineParallelSize                       int      `mapstructure:"pipeline_parallel_size" validate:"omitempty"`
	TensorParallelSize                         int      `mapstructure:"tensor_parallel_size" validate:"omitempty"`
	MaxParallelLoadingWorkers                  int      `mapstructure:"max_parallel_loading_workers" validate:"omitempty"`
	RayWorkersUseNsight                        bool     `mapstructure:"ray_workers_use_nsight"`
	BlockSize                                  int      `mapstructure:"block_size" validate:"omitempty,oneof=8 16 32 64 128"`
	EnablePrefixCaching                        bool     `mapstructure:"enable_prefix_caching"`
	DisableSlidingWindow                       bool     `mapstructure:"disable_sliding_window"`
	UseV2BlockManager                          bool     `mapstructure:"use_v2_block_manager"`
	NumLookaheadSlots                          int      `mapstructure:"num_lookahead_slots" validate:"omitempty"`
	Seed                                       int      `mapstructure:"seed" validate:"omitempty"`
	SwapSpace                                  int      `mapstructure:"swap_space" validate:"omitempty"`
	CPUOffloadGB                               int      `mapstructure:"cpu_offload_gb" validate:"omitempty"`
	GPUMemoryUtilization                       int      `mapstructure:"gpu_memory_utilization" validate:"omitempty"`
	NumGPUBlocksOverride                       int      `mapstructure:"num_gpu_blocks_override" validate:"omitempty"`
	MaxNumBatchedTokens                        int      `mapstructure:"max_num_batched_tokens" validate:"omitempty"`
	MaxNumSeqs                                 int      `mapstructure:"max_num_seqs" validate:"omitempty"`
	MaxLogprobs                                int      `mapstructure:"max_logprobs" validate:"omitempty"`
	DisableLogStats                            bool     `mapstructure:"disable_log_stats"`
	Quantization                               string   `mapstructure:"quantization" validate:"omitempty,oneof=aqlm awq deepspeedfp tpu_int8 fp8 fbgemm_fp8 modelopt marlin gguf gptq_marlin_24 gptq_marlin awq_marlin gptq compressed-tensors bitsandbytes qqq hqq experts_int8 neuron_quant ipex quark moe_wna16 None"`
	RopeScaling                                string   `mapstructure:"rope_scaling" validate:"omitempty"`
	RopeTheta                                  string   `mapstructure:"rope_theta" validate:"omitempty"`
	HFOverrides                                string   `mapstructure:"hf_overrides" validate:"omitempty"`
	EnforceEager                               bool     `mapstructure:"enforce_eager"`
	MaxSeqLenToCapture                         int      `mapstructure:"max_seq_len_to_capture" validate:"omitempty"`
	DisableCustomAllReduce                     bool     `mapstructure:"disable_custom_all_reduce"`
	TokenizerPoolSize                          int      `mapstructure:"tokenizer_pool_size" validate:"omitempty"`
	TokenizerPoolType                          string   `mapstructure:"tokenizer_pool_type" validate:"omitempty"`
	TokenizerPoolExtraConfig                   string   `mapstructure:"tokenizer_pool_extra_config" validate:"omitempty"`
	LimitMMPerPrompt                           int      `mapstructure:"limit_mm_per_prompt" validate:"omitempty"`
	MMProcessorKwargs                          string   `mapstructure:"mm_processor_kwargs" validate:"omitempty"`
	DisableMMPreprocessorCache                 bool     `mapstructure:"disable_mm_preprocessor_cache"`
	EnableLora                                 bool     `mapstructure:"enable_lora"`
	EnableLoraBias                             bool     `mapstructure:"enable_lora_bias"`
	MaxLoras                                   int      `mapstructure:"max_loras" validate:"omitempty"`
	MaxLoraRank                                int      `mapstructure:"max_lora_rank" validate:"omitempty"`
	LoraExtraVocabSize                         int      `mapstructure:"lora_extra_vocab_size" validate:"omitempty"`
	LoraDtype                                  string   `mapstructure:"lora_dtype" validate:"omitempty,oneof=auto float16 bfloat16"`
	LongLoraScalingFactors                     string   `mapstructure:"long_lora_scaling_factors" validate:"omitempty"`
	MaxCPULoras                                int      `mapstructure:"max_cpu_loras" validate:"omitempty"`
	FullyShardedLoras                          bool     `mapstructure:"fully_sharded_loras"`
	EnablePromptAdapter                        bool     `mapstructure:"enable_prompt_adapter"`
	MaxPromptAdapters                          int      `mapstructure:"max_prompt_adapters" validate:"omitempty"`
	MaxPromptAdapterToken                      int      `mapstructure:"max_prompt_adapter_token" validate:"omitempty"`
	Device                                     string   `mapstructure:"device" validate:"omitempty,oneof=auto cuda neuron cpu openvino tpu xpu hpu"`
	NumSchedulerSteps                          int      `mapstructure:"num_scheduler_steps" validate:"omitempty"`
	MultiStepStreamOutputs                     bool     `mapstructure:"multi_step_stream_outputs"`
	SchedulerDelayFactor                       int      `mapstructure:"scheduler_delay_factor" validate:"omitempty"`
	EnableChunkedPrefill                       bool     `mapstructure:"enable_chunked_prefill"`
	SpeculativeModel                           string   `mapstructure:"speculative_model" validate:"omitempty"`
	SpeculativeModelQuantization               string   `mapstructure:"speculative_model_quantization" validate:"omitempty,oneof=aqlm awq deepspeedfp tpu_int8 fp8 fbgemm_fp8 modelopt marlin gguf gptq_marlin_24 gptq_marlin awq_marlin gptq compressed-tensors bitsandbytes qqq hqq experts_int8 neuron_quant ipex quark moe_wna16 None"`
	NumSpeculativeTokens                       int      `mapstructure:"num_speculative_tokens" validate:"omitempty"`
	SpeculativeDisableMQAScorer                bool     `mapstructure:"speculative_disable_mqa_scorer"`
	SpeculativeDraftTensorParallelSize         int      `mapstructure:"speculative_draft_tensor_parallel_size" validate:"omitempty"`
	SpeculativeMaxModelLen                     int      `mapstructure:"speculative_max_model_len" validate:"omitempty"`
	SpeculativeDisableByBatchSize              bool     `mapstructure:"speculative_disable_by_batch_size"`
	NgramPromptLookupMax                       int      `mapstructure:"ngram_prompt_lookup_max" validate:"omitempty"`
	NgramPromptLookupMin                       int      `mapstructure:"ngram_prompt_lookup_min" validate:"omitempty"`
	SpecDecodingAcceptanceMethod               string   `mapstructure:"spec_decoding_acceptance_method" validate:"omitempty,oneof=rejection_sampler typical_acceptance_sampler"`
	TypicalAcceptanceSamplerPosteriorThreshold float64  `mapstructure:"typical_acceptance_sampler_posterior_threshold" validate:"omitempty"`
	TypicalAcceptanceSamplerPosteriorAlpha     float64  `mapstructure:"typical_acceptance_sampler_posterior_alpha" validate:"omitempty"`
	DisableLogprobsDuringSpecDecoding          bool     `mapstructure:"disable_logprobs_during_spec_decoding"`
	ModelLoaderExtraConfig                     string   `mapstructure:"model_loader_extra_config" validate:"omitempty"`
	IgnorePatterns                             string   `mapstructure:"ignore_patterns" validate:"omitempty"`
	PreemptionMode                             string   `mapstructure:"preemption_mode" validate:"omitempty"`
	ServedModelName                            []string `mapstructure:"served_model_name" validate:"omitempty"`
	QLoraAdapterNameOrPath                     string   `mapstructure:"qlora_adapter_name_or_path" validate:"omitempty"`
	OtlpTracesEndpoint                         string   `mapstructure:"otlp_traces_endpoint" validate:"omitempty"`
	CollectDetailedTraces                      bool     `mapstructure:"collect_detailed_traces"`
	DisableAsyncOutputProc                     bool     `mapstructure:"disable_async_output_proc"`
	SchedulingPolicy                           string   `mapstructure:"scheduling_policy" validate:"omitempty,oneof=fcfs priority"`
	OverrideNeuronConfig                       string   `mapstructure:"override_neuron_config" validate:"omitempty"`
	OverridePoolerConfig                       string   `mapstructure:"override_pooler_config" validate:"omitempty"`
	CompilationConfig                          string   `mapstructure:"compilation_config" validate:"omitempty"`
	KVTransferConfig                           string   `mapstructure:"kv_transfer_config" validate:"omitempty"`
	WorkerCls                                  string   `mapstructure:"worker_cls" validate:"omitempty"`
	GenerationConfig                           string   `mapstructure:"generation_config" validate:"omitempty"`
	OverrideGenerationConfig                   string   `mapstructure:"override_generation_config" validate:"omitempty"`
	EnableSleepMode                            bool     `mapstructure:"enable_sleep_mode"`
	CalculateKVScales                          bool     `mapstructure:"calculate_kv_scales"`
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
	HealthCheck string `mapstructure:"health_check" validate:"required"`
}
