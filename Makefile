VERSION ?= $(shell type VERSION 2>nul || cat VERSION 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: dev build frontend build-all test tidy release release-amd64 release-arm64

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

release-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/zpanel-linux-amd64 ./cmd/zpanel

release-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/zpanel-linux-arm64 ./cmd/zpanel

release: build-all release-amd64 release-arm64
	@cd bin && tar czf zpanel-linux-amd64.tar.gz zpanel-linux-amd64
	@cd bin && tar czf zpanel-linux-arm64.tar.gz zpanel-linux-arm64
	@echo "Release packages in bin/"
