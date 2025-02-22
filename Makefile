.DEFAULT_GOAL := build

.PHONY: build_for_alpine fmt vet build run

build_for_alpine: lint
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo ./cmd/medguardian/

build_for_scratch:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/medguardian/

image:
	docker build -t medguardian:test .

run_container:
	docker run --rm -d --env-file .env.local medguardian:test

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build .

lint: vet
	golangci-lint run ./...

clean:
	go clean

run: vet
	DEV=true go run cmd/medguardian/main.go
