package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := connect.Run(context.Background(), logger); err != nil {
		logger.With("err", err).Error("faild to connect run")
		os.Exit(1)
	}
}
