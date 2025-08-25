package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/spf13/viper"
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
