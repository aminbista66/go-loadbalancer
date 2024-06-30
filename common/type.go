package common

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func MimicBrowser(req *http.Request, s *Server) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Origin", "https://www.youtube.com")
    req.Header.Set("Referer", "https://www.youtube.com/")
	req.Host = s.URL.Host
}

func AddCORSHeaders(resp *http.Response) {
	resp.Header.Set("Access-Control-Allow-Origin", "*")
	resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func HandleAntiProxyMechanism(resp *http.Response) error {
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		resp.Header.Set("Content-Type", "text/html; charset=UTF-8")
	}
	return nil
}

func AddCookies(resp *http.Response) {
	resp.Header.Add("Set-Cookie", "AK_SERVER_TIME=1719734195; expires=Sun, 30-Jun-2024 07:57:15 GMT; path=/; secure")
	resp.Header.Add("Set-Cookie", "geo=GB,,LONDON,51.51,-0.13,2643743; expires=Sun, 30-Jun-2024 07:57:35 GMT; secure")
}

func (s *Server) AdvanceProxy() *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(s.URL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Modify the request headers to mimic a real browser
		MimicBrowser(req, s)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// Handle YouTube's anti-proxy mechanisms
		err := HandleAntiProxyMechanism(resp)
		if err != nil {
			return err
		}

        // Add CORS headers to the response
        AddCORSHeaders(resp)
        // Add cookies to the response
        AddCookies(resp)
		return nil
	}

	return proxy
}

type Config struct {
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []string `json:"servers"`
	ListenPort          string   `json:"listenPort"`
	HealthCheckEndpoint string   `json:"healthCheckEndpoint"`
}
