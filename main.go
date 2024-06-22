package main

import (
	"flag"
	"log/slog"
	"os"

	"lb/internal/balancer"
	"lb/pkg/config"
)

func main() {
	configPath := flag.String("configPath", "./config.yaml", "path to the config file")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	serverConfig := config.GetBalancerConfig()

	logLevel := slog.LevelInfo

	if serverConfig.Debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))

	slog.SetDefault(logger)

	balancer.Run()
}
