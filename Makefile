.PHONY: web build run dev test

web:
	cd web && npm install && npm run build

build: web
	go build -o portfolio ./server

# ADMIN_PASSWORD обязателен; SESSION_SECRET рекомендуется
run: build
	./portfolio

# Терминал 1: ADMIN_PASSWORD=dev go run ./server   Терминал 2: make dev
dev:
	cd web && npm run dev

test:
	go vet ./... && go test ./...
