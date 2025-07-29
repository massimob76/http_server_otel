package main

import (
	"log/slog"
	"os"
)

var log = logger()

func logger() *slog.Logger {
	var level slog.Level
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level.UnmarshalText([]byte(lvl))
	} else {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: level,
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	log.Info("logging level set", "level", level.String())
	return log
}
