package common

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL               *url.URL
	ActiveConnections int
	Mutex             sync.Mutex
	Healthy           bool
}

func (s *Server) Proxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

type Config struct {
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []string `json:"servers"`
	ListenPort          string   `json:"listenPort"`
}
