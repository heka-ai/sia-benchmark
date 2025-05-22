package injection

import (
	apiConfig "github.com/heka-ai/benchmark-api/internal/config"
	"github.com/heka-ai/benchmark-api/internal/log"
	api_http "github.com/heka-ai/benchmark-api/internal/web"
	"github.com/heka-ai/benchmark-api/pkg/benchmark"
	"github.com/heka-ai/benchmark-api/pkg/vllm"
	"github.com/ipfans/fxlogger"
	"go.uber.org/fx"
)

var logger = log.GetLogger("injection")

type AppInjector struct {
	Injector *fx.App
}

func NewAppInjector() *AppInjector {
	logger.Info().Msg("Starting the dependency injection")

	app := fx.New(
		fx.WithLogger(fxlogger.WithZerolog(log.GetMainLogger())),

		apiConfig.ConfigFX,
		api_http.HttpModule,
		vllm.VLLMModule,
		benchmark.BenchmarkModule,

		fx.Invoke(func(s *api_http.HttpServer) {}),
	)

	return &AppInjector{
		Injector: app,
	}
}

func (i *AppInjector) Start() {
	logger.Info().Msg("Starting the dependency injection")

	i.Injector.Run()
}
