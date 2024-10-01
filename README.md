# JWT Web Server

A simple Go-based web server and client application implementing JWT authentication using Gorilla Mux and `golang-jwt/jwt`.

## Prerequisites

- **Go**: Version 1.22 or later. [Download Go](https://golang.org/dl/)
- **Docker**: Required for linting with golangci-lint. [Download Docker](https://www.docker.com/get-started)
- **Make**

## Installation

1. **Clone the Repository:**

        git clone https://github.com/yourusername/jwt-web-server.git
        cd jwt-web-server

2. **Install Dependencies:**

        make update_deps

## Usage

### Build

To build both the server and client binaries:

        make build

This will compile the binaries and place them in the `bin/` directory along with the `secret.key`.

### Run Server

To run the server:

        make run_server

### Run Client

To run the client:

        make run_client

## License

This project is licensed under the BSD 2-Clause License. See the [LICENSE](LICENSE) file for details.
