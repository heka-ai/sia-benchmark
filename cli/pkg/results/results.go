package results

type Results struct {
	Date                 *string      `json:"date"`
	Backend              *string      `json:"backend"`
	ModelID              *string      `json:"model_id"`
	TokenizerID          *string      `json:"tokenizer_id"`
	BestOf               *int         `json:"best_of"`
	NumPrompts           *int         `json:"num_prompts"`
	Input                *[]string    `json:"input"`
	ExpectedOutput       *[]string    `json:"expected_output"`
	RequestRate          *string      `json:"request_rate"`
	Duration             *float64     `json:"duration"`
	Completed            *int         `json:"completed"`
	TotalInputTokens     *int         `json:"total_input_tokens"`
	TotalOutputTokens    *int         `json:"total_output_tokens"`
	RequestThroughput    *float64     `json:"request_throughput"`
	OutputThroughput     *float64     `json:"output_throughput"`
	TotalTokenThroughput *float64     `json:"total_token_throughput"`
	InputLens            *[]int       `json:"input_lens"`
	OutputLens           *[]int       `json:"output_lens"`
	Ttfts                *[]float64   `json:"ttfts"`
	Itls                 *[][]float64 `json:"itls"`
	GeneratedTexts       *[]string    `json:"generated_texts"`
	Errors               *[]string    `json:"errors"`
	MeanTtftMs           *float64     `json:"mean_ttft_ms"`
	MedianTtftMs         *float64     `json:"median_ttft_ms"`
	StdTtftMs            *float64     `json:"std_ttft_ms"`
	P99TtftMs            *float64     `json:"p99_ttft_ms"`
	MeanTpotMs           *float64     `json:"mean_tpot_ms"`
	MedianTpotMs         *float64     `json:"median_tpot_ms"`
	Results              *Result      `json:"results"`
	Environment          *Environment `json:"environment"`
	Model                *Model       `json:"model"`
	Task                 *Task        `json:"task"`
	Benchmark            *Benchmark   `json:"benchmark"`
	Dataset              *Dataset     `json:"dataset"`
	Evaluation           *Evaluation  `json:"evaluation"`
	BenchmarkID          *string      `json:"benchmark_id"`
	DatasetPath          *string      `json:"dataset_path"`
	DatasetRevision      *string      `json:"dataset_revision"`
	DatasetSplit         *string      `json:"dataset_split"`
}

type Environment struct {
	Id                 *string `json:"id"`
	Regions            *string `json:"regions"`
	Ec2CpuInstanceType *string `json:"ec2_cpu_instance_type"`
	Ec2GpuInstanceType *string `json:"ec2_gpu_instance_type"`
}

type Model struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type Task struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type Benchmark struct {
	Id             *string `json:"id"`
	Fk_model       *string `json:"fk_model"`
	Fk_environment *string `json:"fk_environment"`
	Fk_task        *string `json:"fk_task"`
	Fk_dataset     *string `json:"fk_dataset"`
	Date           *string `json:"date"`
}

type Dataset struct {
	Id       *string `json:"id"`
	Url      *string `json:"url"`
	Revision *string `json:"revision"`
	Split    *string `json:"split"`
}

type Evaluation struct {
	Id                  *string             `json:"id"`
	EvaluationModel     *string             `json:"evaluation_model"`
	PromptTemplate      *string             `json:"prompt_template"`
	TopK                *int                `json:"top_k"`
	ShowIndicator       *bool               `json:"show_indicator"`
	PrintResults        *bool               `json:"print_results"`
	WriteCache          *bool               `json:"write_cache"`
	UseCache            *bool               `json:"use_cache"`
	SkipOnMissingParams *bool               `json:"skip_on_missing_params"`
	VerboseMode         *bool               `json:"verbose_mode"`
	ThrottleValue       *int                `json:"throttle_value"`
	MetricsDesired      *map[string]float64 `json:"metrics_desired"`
}

type Result struct {
	Id             *[]string    `json:"id"`
	Input          *[]string    `json:"input"`
	ExpectedOutput *[]string    `json:"expected_output"`
	ActualOutput   *[]string    `json:"actual_output"`
	Itls           *[][]float64 `json:"itls"`
	Ttfts          *[]float64   `json:"ttfts"`
}
