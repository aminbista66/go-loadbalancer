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

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s\n", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s\n", err.Error())
	}

	var servers []*common.Server
	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		servers = append(servers, &common.Server{URL: u})
	}

	for _, server := range servers {
		go func(s *common.Server) {
			for range time.Tick(healthCheckInterval) {
				res, err := http.Get(s.URL.String() + config.HealthCheckEndpoint)
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		server := algo.NextServerLeastActive(servers)
		log.Println("Proxying request to server", server.URL.String())
		server.Mutex.Lock()
		server.ActiveConnections++
		server.Mutex.Unlock()
		server.Proxy().ServeHTTP(rw, req)
		server.Mutex.Lock()
		server.ActiveConnections--
		server.Mutex.Unlock()
	})

	log.Println("Starting server on port", config.ListenPort)
	hostUrl := fmt.Sprintf("localhost:%v", config.ListenPort)
	err = http.ListenAndServe(hostUrl, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}

}