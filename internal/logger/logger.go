package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Logger struct {
	log *slog.Logger
}

func (logger *Logger) Info(msg string, v ...interface{}) {
	logger.log.Info(msg, v...)
}

func (logger *Logger) Error(err error) {
	if err != nil {
		logger.log.Error(err.Error())
	}
}

func (logger *Logger) ErrorWithMsg(msg string, err error) {
	logger.log.Error(msg, slog.String("error", err.Error()))
}

func NewLogger(env string) *Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &Logger{log: log}
}
