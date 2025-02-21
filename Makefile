.DEFAULT_GOAL := build

.PHONY: build_for_alpine fmt vet build run

build_for_alpine:
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo .

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build .

clean:
	go clean

run: vet
	DEV=true go run main.go
