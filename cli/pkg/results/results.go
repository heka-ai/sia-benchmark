package results

type Results struct {
	Results          *Result           `json:"results"`
	Environment      *Environment      `json:"environment"`
	Model            *Model            `json:"model"`
	Task             *Task             `json:"task"`
	Benchmark        *Benchmark        `json:"benchmark"`
	Dataset          *Dataset          `json:"dataset"`
	Metrics          *Metrics          `json:"metrics"`
	Evaluation       *Evaluation       `json:"evaluation"`
	BenchmarkID      *string           `json:"benchmark_id"`
	InferenceEngine  *string           `json:"inference_engine"`
	DeploymentParams *DeploymentParams `json:"deployment_params"`
}

type Environment struct {
	Provider        *string `json:"provider"`
	CPUInstanceType *string `json:"cpu_instance_type"`
	GPUInstanceType *string `json:"gpu_instance_type"`
	Region          *string `json:"region"`
}

type Model struct {
	Name        *string `json:"name"`
	TokenizerID *string `json:"tokenizer_id"`
	URL         *string `json:"url"`
	BestOf      *int    `json:"best_of"`
}

type Task struct {
	// Task is typically empty in the API response
}

type Benchmark struct {
	Date                 *string             `json:"date"`
	Metrics              *map[string]float64 `json:"metrics"`
	RequestRate          *string             `json:"request_rate"`
	Duration             *float64            `json:"duration"`
	Completed            *int                `json:"completed"`
	TotalOutputTokens    *int                `json:"total_output_tokens"`
	RequestThroughput    *float64            `json:"request_throughput"`
	OutputThroughput     *float64            `json:"output_throughput"`
	TotalTokenThroughput *float64            `json:"total_token_throughput"`
	NumPrompts           *int                `json:"num_prompts"`
}

type Dataset struct {
	URL      *string `json:"url"`
	Revision *string `json:"revision"`
	Split    *string `json:"split"`
	Name     *string `json:"name"`
}

type Metrics struct {
	AnswerRelevancyMetric     *float64 `json:"AnswerRelevancyMetric"`
	BiasMetric                *float64 `json:"BiasMetric"`
	FaithfulnessMetric        *float64 `json:"FaithfulnessMetric"`
	ContextualPrecisionMetric *float64 `json:"ContextualPrecisionMetric"`
	ContextualRecallMetric    *float64 `json:"ContextualRecallMetric"`
	ContextualRelevancyMetric *float64 `json:"ContextualRelevancyMetric"`
}

type Evaluation struct {
	EvaluationModel *string             `json:"evaluation_model"`
	PromptTemplate  *string             `json:"prompt_template"`
	TopK            *int                `json:"top_k"`
	MetricsDesired  *map[string]float64 `json:"metrics_desired"`
}

type DeploymentParams struct {
	Model              *string `json:"model"`
	Dtype              *string `json:"dtype"`
	Seed               *int    `json:"seed"`
	KVCacheDtype       *string `json:"kv-cache-dtype"`
	MaxModelLen        *string `json:"max-model-len"`
	MaxNumSeqs         *int    `json:"max_num_seqs"`
	TokenizerMode      *string `json:"tokenizer_mode"`
	TensorParallelSize *string `json:"tensor_parallel_size"`
}

type Result struct {
	Input          *[]string    `json:"input"`
	ExpectedOutput *[]string    `json:"expected_output"`
	ActualOutput   *[]string    `json:"actual_output"`
	Itls           *[][]float64 `json:"itls"`
	Ttfts          *[]float64   `json:"ttfts"`
	InputLens      *[]int       `json:"input_lens"`
	OutputLens     *[]int       `json:"output_lens"`
}
