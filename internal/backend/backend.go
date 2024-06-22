package backend

import (
	"log/slog"
	"net"
	"time"

	"lb/pkg/config"
)

func New(addr string) *Backend {
	return &Backend{
		Addr:  addr,
		Alive: true,
	}
}
func (b *Backend) IsAlive() bool {
	return b.Alive
}

func (b *Backend) setAlive(alive bool) {
	b.Alive = alive
}

func (b *Backend) GetConnection() (net.Conn, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	conn, err := net.DialTimeout("tcp", b.Addr, time.Duration(config.GetHealth().Timeout)*time.Second)
	if err != nil {
		slog.Error(err.Error(), "Addr", b.Addr)

		return nil, err
	}

	return conn, nil
}
