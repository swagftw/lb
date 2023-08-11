package main

import (
    "log"
    "os"

    "lb/internal/balancer"
    "lb/pkg/config"
)

func main() {
    lb := new(balancer.LB)

    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        configPath = "./pkg/config/config.yaml"
    }

    err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatal(err)
    }

    lb.LoadBackends()
    lb.Start()
}
