build:
	go build -o bin/github.com/bbapp-org/auth-service main.go

run:
	go run main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

migrate:
	go run scripts/migrate.go

seed:
	go run scripts/seed.go

docker-build:
	docker build -t github.com/bbapp-org/auth-service .

docker-run:
	docker run -p 3001:3001 github.com/bbapp-org/auth-service

clean:
	rm -f bin/github.com/bbapp-org/auth-service
	rm -f coverage.out
	rm -f coverage.html

lint:
	golangci-lint run

fmt:
	go fmt ./...

mod-tidy:
	go mod tidy

dev:
	air

.PHONY: build run test test-coverage migrate seed docker-build docker-run clean lint fmt mod-tidy dev
