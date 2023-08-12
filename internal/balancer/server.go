package balancer

import (
    "fmt"
    "log/slog"
    "math"
    "net/http"
    "net/http/httputil"

    "lb/internal/backend"
)

type Server struct {
    Port int
    lb   *LB
}

// Start starts the reverse proxy server.
func (s *Server) Start() {
    reverseProxy := httputil.ReverseProxy{Director: s.Director}

    go func() {
        slog.Info("starting server", "port", s.Port)

        err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), &reverseProxy)
        if err != nil {
            slog.Error(err.Error())
        }
    }()

    select {}
}

// Director is used to modify the request before it is forwarded to the backend.
// it also selects the backend to forward the request to according to the load balancing algorithm.
// load balancing algorithm: least connections.
func (s *Server) Director(req *http.Request) {
    req.URL.Scheme = "http"

    connections := math.Inf(1)
    var chosenBackend *backend.Backend

    for _, be := range s.lb.backends {
        if be.IsAlive() && float64(be.TotalConnections) < connections {
            connections = float64(be.TotalConnections)
            chosenBackend = be
        }
    }

    chosenBackendIP := chosenBackend.IP
    chosenBackend.TotalConnections++

    if chosenBackendIP == "" {
        slog.Error("no backends are available")

        return
    }

    req.URL.Host = chosenBackendIP
}
