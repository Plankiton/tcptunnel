package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Flag int

func setupLogger(ctx context.Context) *zerolog.Event {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Info()

	for name, value := range ctx.Value("logger_fields").(map[string]interface{}) {
		logger = logger.Interface(name, value)
	}

	return logger
}

func Print(ctx context.Context, message string) {
	setupLogger(ctx).
		Msg(message)
}

func Log(ctx context.Context) *zerolog.Event {
	return setupLogger(ctx)
}

const (
	InfoFlag  = 0
	ErrorFlag = 1
)
