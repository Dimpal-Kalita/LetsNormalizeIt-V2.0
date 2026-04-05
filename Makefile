.PHONY: build run test clean tidy setup env-setup docker-up docker-down docker-logs

setup: env-setup tidy
	@echo "Project setup complete!"

env-setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please update .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./... -v -race

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

tidy:
	go mod tidy

dev:
	@echo "Starting development server with hot reload..."
	@mkdir -p tmp
	air

prod: tidy
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/server cmd/server/main.go

docker-up:
	cd deploy/docker && docker-compose up -d
	@echo "Services started. API available at http://localhost:8080"

docker-down:
	cd deploy/docker && docker-compose down

docker-logs:
	cd deploy/docker && docker-compose logs -f

docker-build:
	cd deploy/docker && docker-compose build

docker-clean:
	cd deploy/docker && docker-compose down -v

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
