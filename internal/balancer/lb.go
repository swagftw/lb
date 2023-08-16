package balancer

import (
	"log"
	"log/slog"
	"os"
	"sync"

	"lb/internal/backend"
	"lb/pkg/config"
)

type LB struct {
	backends  []*backend.Backend
	nextIndex uint32
	mutex     sync.Mutex
}

// LoadBackends loads backends from config
func (lb *LB) LoadBackends() {
	// clear backends
	clear(lb.backends)

	backends := config.GetBackends()

	for _, b := range backends {
		lb.backends = append(lb.backends, backend.New(b.Addr))
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
		backends: make([]*backend.Backend, 0),
	}

	// load backends
	lb.LoadBackends()

	// keep pinging backends until context is cancelled,
	// if the context is cancelled start again
	go lb.pingBackends()

	// start server
	server := Server{
		port: config.GetServerConfig().Port,
		lb:   lb,
	}

	server.Start()
}

// pingBackends pings all backends from the config, concurrently.
func (lb *LB) pingBackends() {
	for _, be := range lb.backends {
		go be.Ping()
	}
}

// round-robin load balancing algorithm
func (lb *LB) next() *backend.Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	totalBackends := len(lb.backends)

	if totalBackends == 0 {
		return nil
	}

	be := lb.backends[0]
	if !be.IsAlive() {
		lb.next()
	}

	if totalBackends >= 2 {
		lb.backends = append(lb.backends[1:], be)
	}

	return be
}
