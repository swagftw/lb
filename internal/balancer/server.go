package balancer

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	port int
	lb   *LB
}

// Start starts the reverse proxy server.
func (s *Server) Start() {
	reverseProxy := httputil.ReverseProxy{
		Director: s.director,
		// Transport: &http.Transport{
		//     MaxIdleConns:    100,
		//     IdleConnTimeout: time.Minute,
		// },
	}

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: &reverseProxy,
	}

	go func() {
		slog.Info("starting server", "port", s.port)

		err := svr.ListenAndServe()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	select {}
}

func (s *Server) director(req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "err", r)
		}
	}()

	be := s.lb.next()
	if be == nil {
		slog.Error("no backends available")

		return
	}

	req.URL.Scheme = "http"
	req.URL.Host = be.Addr
}
