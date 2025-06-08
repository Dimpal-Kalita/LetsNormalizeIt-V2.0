.PHONY: build run test clean tidy setup env-setup

# Setup project for first time
setup: env-setup tidy
	@echo "Project setup complete!"

# Setup environment file
env-setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please update .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Default build target
build:
	go build -o bin/server cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Tidy dependencies
tidy:
	go mod tidy

# Build and run
dev: tidy build run

# Build for production
prod: tidy
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/server cmd/server/main.go
