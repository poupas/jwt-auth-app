SERVER_DIR=server
CLIENT_DIR=client
BIN_DIR=bin
SERVER_BINARY=$(BIN_DIR)/server
CLIENT_BINARY=$(BIN_DIR)/client
BINARY_NAME=jwt-auth-app

.PHONY: all build server client test clean update_deps update_dep check_deps lint run_client run_server


all: build $(BIN_DIR)/secret.key

$(BIN_DIR)/secret.key:
	@echo Creating 256-bit JWT secret key...
	@dd if=/dev/urandom of=$(BIN_DIR)/secret.key bs=1 count=32 >/dev/null 2>&1

build: $(SERVER_BINARY) $(CLIENT_BINARY)

$(SERVER_BINARY): ./$(SERVER_DIR)/*.go
	@echo "Building server..."
	mkdir -p $(BIN_DIR)
	go build -o $(SERVER_BINARY) ./$(SERVER_DIR)

$(CLIENT_BINARY): ./$(CLIENT_DIR)/*.go
	@echo "Building client..."
	mkdir -p $(BIN_DIR)
	go build -o $(CLIENT_BINARY) ./$(CLIENT_DIR)

test:
	@echo "Running tests..."
	go test ./middleware/...

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(SERVER_BINARY)
	rm -f $(CLIENT_BINARY)

fmt:
	go fmt ./...

lint: fmt
	@command -v docker >/dev/null || { echo "You need Docker installed to run the linter" && exit 1; }
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint golangci-lint run -v

update_deps:
	@echo "Updating all dependencies to latest versions..."
	go get -u ./...
	go mod tidy

check_deps:
	@echo "Checking for available dependency updates..."
	go list -u -m all

run_client: $(BIN_DIR)/secret.key
	$(CLIENT_BINARY) -secret $(BIN_DIR)/secret.key

run_server: $(BIN_DIR)/secret.key
	$(SERVER_BINARY) -secret $(BIN_DIR)/secret.key


