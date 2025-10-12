.PHONY: run build test clean dev docker-up docker-down

# Development
dev:
    go run main.go

run:
    go run main.go

# Build
build:
    go build -o bin/main main.go

# Testing
test:
    go test -v ./...

# Clean
clean:
    rm -rf bin/

# Docker
docker-up:
    docker-compose up -d

docker-down:
    docker-compose down

docker-build:
    docker build -t myapp .

# Database migrations (if you add a migration tool later)
migrate-up:
    migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
    migrate -path ./migrations -database "$(DATABASE_URL)" down