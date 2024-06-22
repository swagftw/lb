package balancer

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"lb/pkg/config"
)

type LB struct {
	backends   []*Backend
	httpClient *http.Client
	mutex      sync.Mutex
}

type Backend struct {
	Addr  string
	Alive bool
	ID    int
}

// LoadBackends loads backends from config
func (lb *LB) LoadBackends() {
	backends := config.GetBackends()

	if lb.backends == nil {
		lb.backends = make([]*Backend, 0, len(backends))
	}

	clear(lb.backends)

	for _, b := range backends {
		lb.backends = append(lb.backends, &Backend{
			Addr: b.Addr,
			ID:   b.ID,
		})
	}
}

func Run() {
	lb := &LB{
		httpClient: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: time.Duration(config.GetHealth().Timeout) * time.Second,
			},
		},
	}

	// keep pinging backends until context is cancelled,
	// if the context is cancelled start again
	go lb.pingBackends()

	// start server
	server := Server{
		port: config.GetBalancerConfig().Port,
		lb:   lb,
	}

	server.Start()
}

// pingBackends pings all backends from the config, concurrently.
func (lb *LB) pingBackends() {
	lb.LoadBackends()

	ticker := time.NewTicker(time.Duration(config.GetHealth().Interval) * time.Second)

	for {
		// comment for better spacing
		if <-ticker.C; true {
			// comment for better spacing
			for _, be := range lb.backends {
				// comment for better spacing
				go func(be *Backend) {
					defer func() {
						if r := recover(); r != nil {
							slog.Error("panic while pinging backend", "panic", r)
						}
					}()

					lb.Ping(be)
				}(be)
			}
		}
	}
}

func (lb *LB) Ping(be *Backend) {
	if be == nil {
		return
	}

	resp, err := lb.httpClient.Get(be.Addr)
	if err != nil {
		be.Alive = false

		slog.Error("error while sending ping", "err", err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		be.Alive = false

		return
	}

	be.Alive = true
}

// round-robin load balancing algorithm
func (lb *LB) next() *Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	totalBackends := len(lb.backends)

	if totalBackends == 0 {
		return nil
	}

	be := lb.backends[0]

	be.Alive = true
	if !be.Alive {
		lb.next()
	}

	// this is smart, if the first backend is used push to back of the slice
	if totalBackends >= 2 {
		lb.backends = append(lb.backends[1:], be)
	}

	return be
}
