package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"GopherChessParty/internal/interfaces"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Logger struct {
	log *slog.Logger
}

// getCallerInfo возвращает информацию о месте вызова
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 2 - пропускаем getCallerInfo и метод логгера
	if !ok {
		return "unknown"
	}
	// Получаем только имя файла и директории
	dir, file := filepath.Split(file)
	// Берем только последнюю директорию
	dirs := filepath.SplitList(dir)
	if len(dirs) > 0 {
		dir = dirs[len(dirs)-1]
	}
	return fmt.Sprintf("%s/%s:%d", dir, file, line)
}

func (logger *Logger) Info(msg string, v ...interface{}) {
	caller := getCallerInfo()
	logger.log.Info(msg, append([]interface{}{"caller", caller}, v...)...)
}

func (logger *Logger) Error(err error) {
	if err != nil {
		caller := getCallerInfo()
		logger.log.Error(err.Error(), "caller", caller)
	}
}

func (logger *Logger) ErrorWithMsg(msg string, err error) {
	caller := getCallerInfo()
	logger.log.Error(msg,
		slog.String("error", err.Error()),
		slog.String("caller", caller),
	)
}

func NewLogger(env string) interfaces.ILogger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						return slog.Attr{
							Key: slog.SourceKey,
							Value: slog.StringValue(
								fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line),
							),
						}
					}
					return a
				},
			}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						return slog.Attr{
							Key: slog.SourceKey,
							Value: slog.StringValue(
								fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line),
							),
						}
					}
					return a
				},
			}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						return slog.Attr{
							Key: slog.SourceKey,
							Value: slog.StringValue(
								fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line),
							),
						}
					}
					return a
				},
			}),
		)
	}

	return &Logger{log: log}
}
