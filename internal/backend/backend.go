package backend

import (
    "context"
    "log/slog"
    "net"
    "time"

    "lb/pkg/config"
)

type Backend struct {
    IP               string
    TotalConnections int
    Alive            bool
}

func (b *Backend) IsAlive() bool {
    return b.Alive
}

func (b *Backend) setAlive(alive bool) {
    b.Alive = alive
}

func (b *Backend) Ping(ctx context.Context) {
    timeout := config.GetHealth().Timeout
    interval := config.GetHealth().Interval

    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    defer ticker.Stop()

    for {
        slog.Debug("checking health", "IP", b.IP)

        _, err := net.DialTimeout("tcp", b.IP, time.Duration(timeout)*time.Second)
        if err != nil {
            slog.Error(err.Error(), "IP", b.IP)
            b.setAlive(false)
            b.TotalConnections = 0
        } else {
            b.setAlive(true)
        }

        select {
        case <-ctx.Done():
            return
        default:
            <-ticker.C
        }
    }
}
