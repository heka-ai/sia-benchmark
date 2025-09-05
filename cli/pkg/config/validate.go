package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	log "github.com/heka-ai/benchmark-cli/internal/logs"
)

var logger = log.GetLogger("config")
var config Config

// Init the config and validate it
func Init() {
	InitConfig()
}

func InitFlags() {}

func GetConfig() Config {
	return config
}

func InitConfig() {
	localConfig := Config{}

	filename := viper.GetString("config")
	if strings.TrimSpace(filename) == "" {
		filename = "bench.toml"
	}

	// Enforce .toml extension
	if !strings.HasSuffix(strings.ToLower(filename), ".toml") {
		logger.Fatal().Msg("Config file must have .toml extension")
	}

	viper.SetConfigFile(filename)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read config file")
	}

	err = viper.Unmarshal(&localConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	// Normalize config: treat empty strings as "None" for enum-like fields
	// and clear zero values for pointer ints where omitempty is desired.
	normalizeConfig(&localConfig)

	validate := validator.New()
	err = validate.Struct(localConfig)
	if err != nil {
		// TODO: improve errors handling
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error().Msgf("Failed to validate config, expected %s, got %s (tags: %s)", err.Field(), err.Value(), err.Tag())
			// Add explicit guidance for missing HF token
			if err.Field() == "Token" && err.Tag() == "required" {
				logger.Error().Msg("benchmark.token is required. Set it in [benchmark] token = \"...\" or export one of: HF_TOKEN, HUGGINGFACEHUB_API_TOKEN, HF_API_TOKEN")
			}
			fmt.Println(err)
		}
		os.Exit(1)
	}

	// Provider-specific resource validation
	validateResourcesAllocation(&localConfig)

	logger.Info().Msgf("Config validated successfully")

	config = localConfig
}

func normalizeConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	// Fill benchmark token from common HF env vars if empty
	if cfg.BenchmarkConfig != nil {
		if strings.TrimSpace(cfg.BenchmarkConfig.Token) == "" {
			for _, envKey := range []string{"HF_TOKEN", "HUGGINGFACEHUB_API_TOKEN", "HF_API_TOKEN"} {
				if env := os.Getenv(envKey); strings.TrimSpace(env) != "" {
					cfg.BenchmarkConfig.Token = env
					break
				}
			}
		}
	}

	if cfg.VLLMConfig == nil {
		return
	}

	v := cfg.VLLMConfig
	rv := reflect.ValueOf(v).Elem()
	rt := reflect.TypeOf(*v)

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		// Only handle pointer fields
		if f.Kind() != reflect.Ptr {
			continue
		}
		if f.IsNil() {
			continue
		}
		fv := f.Elem()
		switch fv.Kind() {
		case reflect.String:
			if strings.TrimSpace(fv.String()) == "" {
				val := "None"
				f.Set(reflect.ValueOf(&val))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fv.Int() == 0 {
				f.Set(reflect.Zero(f.Type()))
			}
		case reflect.Float32, reflect.Float64:
			if fv.Float() == 0 {
				f.Set(reflect.Zero(f.Type()))
			}
		case reflect.Bool:
			// Only treat as unset when the raw TOML value was an empty string
			fieldTag := rt.Field(i).Tag.Get("mapstructure")
			if strings.TrimSpace(fieldTag) != "" {
				raw := viper.Get("vllm." + fieldTag)
				if s, ok := raw.(string); ok && strings.TrimSpace(s) == "" {
					f.Set(reflect.Zero(f.Type()))
				}
			}
		default:
			// Leave other types untouched to preserve explicit values
		}
	}
}

func validateResourcesAllocation(cfg *Config) {
	if cfg == nil {
		return
	}

	if cfg.InferenceEngine != "vllm" {
		logger.Info().
			Str("provider", cfg.Provider).
			Str("engine", cfg.InferenceEngine).
			Msg("Skipping GPU memory validation: not vLLM")
		return
	}

	logger.Info().
		Str("provider", cfg.Provider).
		Str("engine", cfg.InferenceEngine).
		Msg("Validating GPU memory against model size requirement")

	// First, try to extract model size in B from the model string. If we can't, skip.
	model := strings.TrimSpace(cfg.VLLMConfig.Model)
	if model == "" {
		logger.Warn().Msg("Skipping GPU memory validation: empty vLLM model name")
		return
	}

	modelB, ok := extractModelSizeBillions(model)
	if !ok {
		logger.Warn().Str("model", model).Msg("Could not infer model size in billions (XB). Skipping GPU memory check")
		return
	}
	requiredGiB := float64(modelB*2 + 5)

	// Fetch GPU total GiB for the current provider
	gpuGiB, okGpu := getGpuTotalGiB(cfg)
	if !okGpu {
		logger.Warn().Msg("Skipping GPU memory validation: could not determine GPU memory for provider")
		return
	}
	logger.Info().Float64("gpu_memory_gib", gpuGiB).Msg("Determined GPU total memory")

	if gpuGiB <= requiredGiB {
		logger.Fatal().
			Str("provider", cfg.Provider).
			Float64("gpu_gib", gpuGiB).
			Int("model_b", modelB).
			Float64("required_gib", requiredGiB).
			Msg("GPU memory is insufficient for the selected model: require > model_size*2 + 5 GiB")
	}

	logger.Info().
		Float64("gpu_gib", gpuGiB).
		Float64("minimum_required_gib", requiredGiB).
		Msg("GPU memory seem to be sufficient for the selected model")
}

// extractModelSizeBillions parses the last occurrence of an "XB" token from a model name.
func extractModelSizeBillions(model string) (int, bool) {
	pattern := regexp.MustCompile(`(?i)(\d+)\s*b`)
	matches := pattern.FindAllStringSubmatch(model, -1)
	if len(matches) == 0 {
		return 0, false
	}
	last := matches[len(matches)-1]
	val, err := strconv.Atoi(last[1])
	if err != nil {
		return 0, false
	}
	return val, true
}

// getGpuTotalGiB determines the total GPU memory in GiB for the configured provider.
// For AWS, it uses the SDK DescribeInstanceTypes; for local, it queries nvidia-smi
// (falling back to tegrastats). Returns (value, true) on success; (0, false) otherwise.
func getGpuTotalGiB(cfg *Config) (float64, bool) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "aws":
		instanceType := strings.TrimSpace(cfg.AWSConfig.GPUInstanceType)
		region := strings.TrimSpace(cfg.AWSConfig.Region)
		if instanceType == "" || region == "" {
			return 0, false
		}

		var awscfg aws.Config
		var err error
		if p := strings.TrimSpace(cfg.AWSConfig.ProfileName); p != "" {
			awscfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
				awsconfig.WithSharedConfigProfile(p),
				awsconfig.WithRegion(region),
			)
		} else {
			awscfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
				awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSConfig.AWSAccessKey, cfg.AWSConfig.AWSSecretKey, "")),
				awsconfig.WithRegion(region),
			)
		}
		if err != nil {
			return 0, false
		}

		ec2client := ec2.NewFromConfig(awscfg)
		resp, err := ec2client.DescribeInstanceTypes(context.TODO(), &ec2.DescribeInstanceTypesInput{
			InstanceTypes: []types.InstanceType{types.InstanceType(instanceType)},
		})
		if err != nil {
			return 0, false
		}
		if len(resp.InstanceTypes) == 0 || resp.InstanceTypes[0].GpuInfo == nil || resp.InstanceTypes[0].GpuInfo.TotalGpuMemoryInMiB == nil {
			return 0, false
		}
		totalMiB := float64(*resp.InstanceTypes[0].GpuInfo.TotalGpuMemoryInMiB)
		return totalMiB / 1024.0, true

	case "local":
		cmd := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader,nounits")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			var totalMiB float64
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l == "" {
					continue
				}
				if v, perr := strconv.Atoi(l); perr == nil {
					totalMiB += float64(v)
				}
			}
			if totalMiB > 0 {
				return totalMiB / 1024.0, true
			}
		}

		if _, lookErr := exec.LookPath("tegrastats"); lookErr == nil {
			cmd = exec.Command("sh", "-c", "tegrastats --interval 200 | head -n 1")
			if tout, terr := cmd.Output(); terr == nil {
				line := strings.TrimSpace(string(tout))
				re := regexp.MustCompile(`RAM\s+\d+\/(\d+)MB`)
				if m := re.FindStringSubmatch(line); len(m) == 2 {
					if mb, conv := strconv.Atoi(m[1]); conv == nil && mb > 0 {
						return float64(mb) / 1024.0, true
					}
				}
			}
		}

		return 0, false
	default:
		return 0, false
	}
}
