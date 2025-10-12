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

migrate-up:
	go run cmd/migrate/main.go -dir ./internal/storage/sql/migrations up 

migrate-down:
	go run cmd/migrate/main.go -dir internal/storage/sql/migrations down 

migrate-status:
	go run cmd/migrate/main.go -dir internal/storage/sql/migrations status

migrate-version:
	go run cmd/migrate/main.go -dir internal/storage/sql/migrations version

migrate-reset:
	go run cmd/migrate/main.go -dir internal/storage/sql/migrations reset

migrate-create:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go -dir internal/storage/sql/migrations create $$name 


sqlc-generate:
	sqlc generate

sqlc-verify:
	sqlc verify
