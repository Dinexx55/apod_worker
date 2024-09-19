# Build the project
build:
	go build -o bin/apod_service ./cmd/main.go

# Run tests
test:
	go test ./internal/service... -v

# Run the service
run:
	go run ./cmd/main.go

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose up with rebuild
docker-up-b:
	docker-compose up --build

# Docker compose down
docker-down:
	docker-compose down

# Docker compose down and clear volumes
docker-down-v:
	docker-compose down -v

# Help information
help:
	@echo "Makefile commands:"
	@echo "  build         Build the Go binary"
	@echo "  test          Run tests"
	@echo "  run           Run the service"
	@echo "  docker-up     Start docker containers"
	@echo "  docker-up-b    Start docker containers with rebuild"
	@echo "  docker-down   Stop docker containers"
	@echo "  docker-down-v   Stop docker containers and clear the volumes"

.PHONY: build test run docker-up docker-up-b docker-down docker-down-v help