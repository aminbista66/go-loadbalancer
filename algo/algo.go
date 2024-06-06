package algo

import (
	"lb/common"
)

func NextServerLeastActive(servers []*common.Server) *common.Server {
	leastActiveConnections := -1
	leastActiveServer := servers[0]

	for _, server := range servers {
		server.Mutex.Lock()
		if (leastActiveConnections == -1 || server.ActiveConnections < leastActiveConnections) && server.Healthy {
			leastActiveConnections = server.ActiveConnections
			leastActiveServer = server
		}
		server.Mutex.Unlock()
	}
	return leastActiveServer
}