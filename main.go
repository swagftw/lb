package main

import (
    "log"
    "log/slog"
    "os"

    "lb/internal/balancer"
    "lb/pkg/config"
)

func main() {
    logLevel := new(slog.LevelVar)
    logLevel.Set(slog.LevelDebug)

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     logLevel,
        AddSource: true,
    }))

    slog.SetDefault(logger)

    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        configPath = "./pkg/config/config.yaml"
    }

    err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatal(err)
    }

    balancer.Run()
}
