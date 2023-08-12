package balancer

import (
    "context"
    "log"
    "log/slog"
    "os"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"

    "lb/internal/backend"
    "lb/pkg/config"
)

type LB struct {
    Ctx      context.Context
    Cancel   context.CancelFunc
    backends map[string]*backend.Backend
}

// LoadBackends loads backends from config
func (lb *LB) LoadBackends() {
    // clear backends
    clear(lb.backends)

    backends := config.GetBackends()

    for _, b := range backends {
        // check if backend already exists
        if _, ok := lb.backends[b.IP]; ok {
            continue
        }

        lb.backends[b.IP] = &backend.Backend{
            IP:               b.IP,
            TotalConnections: 0,
            Alive:            true,
        }
    }
}

func Run() {
    logLevel := new(slog.LevelVar)
    logLevel.Set(slog.LevelInfo)

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

    lb := &LB{
        backends: make(map[string]*backend.Backend),
    }

    // load backends
    lb.LoadBackends()

    lb.Ctx, lb.Cancel = context.WithCancel(context.Background())

    viper.OnConfigChange(func(e fsnotify.Event) {
        slog.Debug("config file changed:", "event", e.Name)
        err = viper.ReadInConfig()
        if err != nil {
            slog.Error(err.Error(), "msg", "error reading config file")

            return
        }

        // load updated backends
        lb.LoadBackends()

        // cancel current context
        lb.Cancel()
    })

    // keep pinging backends until context is cancelled,
    // if the context is cancelled start again
    go func() {
        for {
            lb.pingBackends()
            <-lb.Ctx.Done()

            // once the context is cancelled (config changed), create a new context
            lb.Ctx, lb.Cancel = context.WithCancel(context.Background())
        }
    }()

    // start server
    server := Server{
        Port: config.GetServerConfig().Port,
        lb:   lb,
    }

    server.Start()
}

// pingBackends pings all backends from the config, concurrently.
func (lb *LB) pingBackends() {
    for _, be := range lb.backends {
        go be.Ping(lb.Ctx)
    }
}
