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

# Vue 构建产物复制到 embed 目录（go:embed 不允许 ../ 路径）
sync-frontend:
	cd web && npm run build
	rm -rf internal/web/dist && cp -r web/dist internal/web/dist

test:
	go test ./...

tidy:
	go mod tidy
