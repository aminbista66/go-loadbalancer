# Load Balancer

This project is a simple load balancer written in Go. It proxies HTTP requests to a pool of backend servers, distributing the load based on the least active connections. The load balancer also performs periodic health checks on the backend servers to ensure they are available before routing traffic to them.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Code Overview](#code-overview)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/load-balancer.git
    cd load-balancer
    ```

2. Install the dependencies:
    ```sh
    go mod tidy
    ```

3. Build the application:
    ```sh
    go build -o load-balancer
    ```

## Configuration

The load balancer is configured using a JSON configuration file. Below is an example `config.json`:

```json
{
    "healthCheckInterval": "10s",
    "servers": [
        "http://localhost:8081",
        "http://localhost:8082"
    ],
    "listenPort": "8080",
    "healthCheckEndpoint": "/health"
}
```

- `healthCheckInterval`: Interval between health checks (e.g., "10s" for 10 seconds).
- `servers`: List of backend server URLs.
- `listenPort`: Port on which the load balancer listens.
- `healthCheckEndpoint`: Endpoint used to check the health of the backend servers.

## Usage

1. Ensure the configuration file (`config.json`) is in the same directory as the executable or specify its path when loading the configuration.
2. Start the load balancer:
    ```sh
    ./load-balancer
    ```
3. The load balancer will start listening on the specified port and proxying requests to the backend servers.

## Code Overview

### `main.go`

- `createServerObject(config *common.Config) []*common.Server`: Creates a slice of `Server` objects from the configuration.
- `performHealthCheck(servers []*common.Server, healthCheckEndpoint string, healthCheckInterval time.Duration)`: Performs health checks on the backend servers at the specified interval.
- `main()`: Entry point of the application. It loads the configuration, sets up the server objects, initiates health checks, and starts the HTTP server to proxy requests.

### `config/config.go`

- `LoadConfig(file string) (common.Config, error)`: Loads the configuration from a JSON file.

### `common/common.go`

- `Server`: Struct representing a backend server.
- `Proxy() *httputil.ReverseProxy`: Returns a reverse proxy for the server.
- `Config`: Struct representing the configuration format.

### `algo/algo.go`

- `NextServerLeastActive(servers []*common.Server) *common.Server`: Algorithm to select the next server with the least active connections that is also healthy.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss changes or feature requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.