package balancer

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	port int
	lb   *LB
}

// Start starts the reverse proxy server.
func (s *Server) Start() {
	reverseProxy := &httputil.ReverseProxy{
		Director: s.director,
	}

	errChan := make(chan error, 1)

	go func() {
		slog.Info("starting server", "port", s.port)

		err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), reverseProxy)
		if err != nil {
			slog.Error(err.Error())

			errChan <- err
		}
	}()

	if err := <-errChan; true {
		slog.Error("failed to start load balancer", "err", err)
	}
}

func (s *Server) director(req *http.Request) {
	be := s.lb.next()
	if be == nil {
		slog.Error("no backends available")

		return
	}

	beURL, _ := url.Parse(be.Addr)

	req.URL.Scheme = "http"
	req.URL.Host = beURL.Host
}
