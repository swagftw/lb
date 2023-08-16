package backend

import (
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/pkg/errors"

	"lb/pkg/config"
)

type Backend struct {
	Addr  string
	Alive bool
}

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

func (b *Backend) Ping() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	interval := config.GetHealth().Interval

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		slog.Debug("checking health", "Addr", b.Addr)

		// get connection from pool
		conn, err := b.GetConnection()
		if err != nil {
			// wait for ticker
			<-ticker.C

			continue
		}

		// keep reading from connection
		one := make([]byte, 1)

		// this is a blocking call until timeout hits
		// in this case there is no timeout set, so it will block
		// until connection is closed
		_, err = conn.Read(one)
		if err != nil {
			slog.Error("connection closed", "Addr", b.Addr)
			b.setAlive(false)

			switch {
			case errors.Is(err, net.ErrClosed):
				b.setAlive(false)
				slog.Error("connection closed", "Addr", b.Addr)
			case errors.Is(err, io.EOF):
				err = conn.Close()
				if err != nil {
					slog.Error("failed to close connection", "err", err.Error(), "Addr", b.Addr)
				}

				slog.Error("connection timed out", "Addr", b.Addr)
			default:
				err = conn.Close()
				if err != nil {
					slog.Error("failed to close connection", "err", err.Error(), "Addr", b.Addr)
				}

				slog.Error("error reading from connection", "err", err.Error(), "Addr", b.Addr)
			}
		} else {
			b.setAlive(true)
		}

		// wait for ticker
		<-ticker.C
	}
}
