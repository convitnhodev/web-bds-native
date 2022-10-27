.PHONY: default
default: web;

api:
	@go run cmd/deein/main.go api
.PHONY: api

web:
	@echo "==> web"
	@ulimit -n 2048 && air -c .web.toml
.PHONY: web

build:
	@echo "==> build"
	go build -o ./tmp/main cmd/deein/main.go && chmod +x ./tmp/main
.PHONY: build
