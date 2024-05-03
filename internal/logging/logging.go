package logging

import (
	"log/slog"
	"os"
)

type LogConfiguration struct {
	Level     string
	LogToFile bool
	File      string
}

var logger *slog.Logger

func (lc *LogConfiguration) NewLogger() (*slog.Logger, error) {
	if logger == nil {

		var level slog.Level

		switch lc.Level {
		case "debug", "DEBUG":
			level = slog.LevelDebug
		case "info", "INFO":
			level = slog.LevelInfo
		case "warn", "WARN":
			level = slog.LevelWarn
		case "error", "ERROR":
			level = slog.LevelError
		}
		opts := &slog.HandlerOptions{
			Level: level,
		}

		if lc.LogToFile {
			f, err := os.OpenFile(lc.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				return nil, err
			}
			logger = slog.New(slog.NewTextHandler(f, opts))
		} else {
			logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
		}
	}
	return logger, nil

}
