.PHONY: default
default: api;

api:
	@go run cmd/deein/main.go api
.PHONY: api

web:
	@go run cmd/deein/main.go web
.PHONY: web
