package config

import (
	"flag"
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
	InitFlags()
	InitConfig()
}

func InitFlags() {
	flag.String("config", "bench.toml", "Path to the config file")
	flag.Parse()
}

func GetConfig() Config {
	return config
}

func InitConfig() {
	localConfig := Config{}

	filename := flag.Lookup("config").Value.String()

	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

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

	// Fill benchmark token from HF_TOKEN if empty
	if cfg.BenchmarkConfig != nil {
		if strings.TrimSpace(cfg.BenchmarkConfig.Token) == "" {
			if env := os.Getenv("HF_TOKEN"); strings.TrimSpace(env) != "" {
				cfg.BenchmarkConfig.Token = env
			}
		}
	}

	if cfg.VLLMConfig == nil {
		return
	}

	v := cfg.VLLMConfig
	rv := reflect.ValueOf(v).Elem()

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
		default:
			// Leave bool and other types untouched to preserve explicit values
		}
	}
}
