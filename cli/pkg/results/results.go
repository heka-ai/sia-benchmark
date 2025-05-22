package results

type Results struct {
	Date                 string       `json:"date" validate:"required"`
	Backend              string       `json:"backend" validate:"required"`
	ModelID              string       `json:"model_id" validate:"required"`
	TokenizerID          string       `json:"tokenizer_id" validate:"required"`
	BestOf               int          `json:"best_of" validate:"required"`
	NumPrompts           int          `json:"num_prompts" validate:"required"`
	Input                []string     `json:"input" validate:"required"`
	ExpectedOutput       []string     `json:"expected_output" validate:"required"`
	RequestRate          string       `json:"request_rate" validate:"required"`
	Duration             float64      `json:"duration" validate:"required"`
	Completed            int          `json:"completed" validate:"required"`
	TotalInputTokens     int          `json:"total_input_tokens" validate:"required"`
	TotalOutputTokens    int          `json:"total_output_tokens" validate:"required"`
	RequestThroughput    float64      `json:"request_throughput" validate:"required"`
	OutputThroughput     float64      `json:"output_throughput" validate:"required"`
	TotalTokenThroughput float64      `json:"total_token_throughput" validate:"required"`
	InputLens            []int        `json:"input_lens" validate:"required"`
	OutputLens           []int        `json:"output_lens" validate:"required"`
	Ttfts                []int        `json:"ttfts" validate:"required"`
	Itls                 []int        `json:"itls" validate:"required"`
	GeneratedTexts       []string     `json:"generated_texts" validate:"required"`
	Errors               []string     `json:"errors" validate:"required"`
	MeanTtftMs           float64      `json:"mean_ttft_ms" validate:"required"`
	MedianTtftMs         float64      `json:"median_ttft_ms" validate:"required"`
	StdTtftMs            float64      `json:"std_ttft_ms" validate:"required"`
	P99TtftMs            float64      `json:"p99_ttft_ms" validate:"required"`
	MeanTpotMs           float64      `json:"mean_tpot_ms" validate:"required"`
	MedianTpotMs         float64      `json:"median_tpot_ms" validate:"required"`
	Results              *Result      `json:"results" validate:"required"`
	Environment          *Environment `json:"environment" validate:"required"`
	Model                *Model       `json:"model" validate:"required"`
	Task                 *Task        `json:"task" validate:"required"`
	Benchmark            *Benchmark   `json:"benchmark" validate:"required"`
	Dataset              *Dataset     `json:"dataset" validate:"required"`
	Evaluation           *Evaluation  `json:"evaluation" validate:"required"`
	BenchmarkID          string       `json:"benchmark_id" validate:"required"`
	DatasetPath          string       `json:"dataset_path" validate:"required"`
	DatasetRevision      string       `json:"dataset_revision" validate:"required"`
	DatasetSplit         string       `json:"dataset_split" validate:"required"`
}

type Environment struct {
	Id                 string `json:"id" validate:"required"`
	Regions            string `json:"regions" validate:"required"`
	Ec2CpuInstanceType string `json:"ec2_cpu_instance_type" validate:"required"`
	Ec2GpuInstanceType string `json:"ec2_gpu_instance_type" validate:"required"`
}

type Model struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Task struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Benchmark struct {
	Id             string `json:"id" validate:"required"`
	Fk_model       string `json:"fk_model" validate:"required"`
	Fk_environment string `json:"fk_environment" validate:"required"`
	Fk_task        string `json:"fk_task" validate:"required"`
	Fk_dataset     string `json:"fk_dataset" validate:"required"`
	Date           string `json:"date" validate:"required"`
}

type Dataset struct {
	Id       string `json:"id" validate:"required"`
	Url      string `json:"url" validate:"required"`
	Revision string `json:"revision" validate:"required"`
	Split    string `json:"split" validate:"required"`
}

type Evaluation struct {
	Id                  string             `json:"id" validate:"required"`
	EvaluationModel     string             `json:"evaluation_model" validate:"required"`
	PromptTemplate      string             `json:"prompt_template" validate:"required"`
	TopK                int                `json:"top_k" validate:"required"`
	ShowIndicator       bool               `json:"show_indicator" validate:"required"`
	PrintResults        bool               `json:"print_results" validate:"required"`
	WriteCache          bool               `json:"write_cache" validate:"required"`
	UseCache            bool               `json:"use_cache" validate:"required"`
	SkipOnMissingParams bool               `json:"skip_on_missing_params" validate:"required"`
	VerboseMode         bool               `json:"verbose_mode" validate:"required"`
	ThrottleValue       int                `json:"throttle_value" validate:"required"`
	MetricsDesired      map[string]float64 `json:"metrics_desired" validate:"required"`
}

type Result struct {
	Id             []string    `json:"id" validate:"required"`
	Input          []string    `json:"input" validate:"required"`
	ExpectedOutput []string    `json:"expected_output" validate:"required"`
	ActualOutput   []string    `json:"actual_output" validate:"required"`
	Itls           [][]float64 `json:"itls" validate:"required"`
	Ttfts          [][]float64 `json:"ttfts" validate:"required"`
}
