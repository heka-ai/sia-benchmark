package log

import (
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func GetLogger(module string) zerolog.Logger {
	return GetMainLogger().With().Str("module", module).Timestamp().Logger()
}

func GetMainLogger() zerolog.Logger {
	zlog.Level(zerolog.InfoLevel)
	return zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}
