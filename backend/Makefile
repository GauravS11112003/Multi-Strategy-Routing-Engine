.PHONY: run build clean test lint deps worker

run:
	go run cmd/main.go

worker:
	go run cmd/worker/main.go

build:
	CGO_ENABLED=0 go build -o bin/api cmd/main.go
	CGO_ENABLED=0 go build -o bin/worker cmd/worker/main.go

clean:
	rm -rf bin/

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

deps:
	go mod download
	go mod tidy
