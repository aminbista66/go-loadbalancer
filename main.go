package main

import (
	"fmt"
	"lb/algo"
	"lb/common"
	"lb/config"
	"log"
	"net/http"
	"net/url"
	"time"
)


func createServerObject(config *common.Config) []*common.Server {
	var servers []*common.Server
	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		servers = append(servers, &common.Server{URL: u})
	}
	return servers
}

func performHealthCheck(servers []*common.Server, healthCheckEndpoint string, healthCheckInterval time.Duration) {
	for _, server := range servers {
		go func(s *common.Server) {
			for range time.Tick(healthCheckInterval) {
				res, err := http.Get(s.URL.String() + healthCheckEndpoint)
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}
}


func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s\n", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s\n", err.Error())
	}

	servers := createServerObject(&config)

	performHealthCheck(servers, config.HealthCheckEndpoint, healthCheckInterval)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		server := algo.NextServerLeastActive(servers)
		log.Printf("[ %v ] Request on current server: %v", req.Method, req.URL.String())
		log.Printf("[ %v ] Proxy request to server url: %v \n", req.Method, server.URL.Hostname())
		server.Mutex.Lock()
		server.ActiveConnections++
		server.Mutex.Unlock()
		req.Header.Set("Origin", server.URL.String())
		req.Header.Set("Referer", server.URL.String())
		req.Header.Set("Access-Control-Allow-Origin", "*")
		server.AdvanceProxy().ServeHTTP(rw, req)
		server.Mutex.Lock()
		server.ActiveConnections--
		server.Mutex.Unlock()
	})

	log.Println("Starting server on port", config.ListenPort)
	hostUrl := fmt.Sprintf(":%v", config.ListenPort)
	err = http.ListenAndServe(hostUrl, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}
}