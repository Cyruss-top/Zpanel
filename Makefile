VERSION ?= $(shell type VERSION 2>nul || cat VERSION 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: dev build frontend test tidy

dev:
	go run ./cmd/zpanel server

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/zpanel ./cmd/zpanel

frontend:
	cd web && npm ci && npm run build

build-all: frontend build

test:
	go test ./...

tidy:
	go mod tidy
