package balancer

import (
    "lb/internal/backend"
    "lb/pkg/config"
)

type LB struct {
    backends map[string]*backend.Backend
}

// LoadBackends loads backends from config
func (lb *LB) LoadBackends() {
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

func (lb *LB) Start() {

}

func (lb *LB) pingBackends() {
    // ip is IP address and be is backend
    for ip, be := range lb.backends {
        go func() {
            config.GetHealth().Interval
        }()
    }
}

func (lb *LB) ping() {
    ping
}
