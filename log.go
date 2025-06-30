package main

import (
	"log/slog"
	"os"
	"strings"
)

func setupLogging() {
	var logger *slog.Logger

	if strings.EqualFold(os.Getenv("ENV"), "development") {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	} else {
		writer, err := os.OpenFile("kora.log", os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		logger = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)
}
