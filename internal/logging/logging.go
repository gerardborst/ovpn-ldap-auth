/*
OpenVPN ldap auth - OpenVPN Ldap authentication

Copyright (C) 2019 - 2021 Egbert Pot
Copyright (C) 2021 - 2024 Gerard Borst

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package logging

import (
	"log"
	"log/slog"
	"os"
)

type LogConfiguration struct {
	Level     string
	LogToFile bool
	File      string
}

var logger *slog.Logger

func NewLogger(lc *LogConfiguration) *slog.Logger {
	if logger == nil {
		if lc == nil {
			log.Fatal("log configuration is nil")
		}

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
				log.Fatalf("unable to initialize logger, %v", err)
			}
			logger = slog.New(slog.NewTextHandler(f, opts))
		} else {
			logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
		}
	}
	return logger

}

func GetLogger() *slog.Logger {
	if logger == nil {
		log.Fatal("log configuration is nil")
	}
	return logger
}
