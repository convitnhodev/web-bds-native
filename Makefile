<<<<<<< HEAD
.PHONY: default
default: web;

api:
	@go run cmd/deein/main.go api
.PHONY: api

web:
	@echo "==> web"
	@ulimit -n 2048 && air -c .web.toml
.PHONY: web

css:
	@echo "==> css"
	cd node && pnpm dev
.PHONY: css

build:
	@echo "==> build"
	go build -o ./tmp/main cmd/deein/main.go && chmod +x ./tmp/main
.PHONY: build
=======
project_name = deeincom
image_name = deeincom:latest

run-local:
	./air serve

run-css-dev:
	cd node && pnpm dev

requirements:
	go mod tidy

clean-packages:
	go clean -modcache

up:
	make up-silent
	make shell

build:
	docker build -t $(image_name) .

build-no-cache:
	docker build --no-cache -t $(image_name) .

up-silent:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app

up-silent-prefork:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app -prod

delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

shell:
	docker exec -it $(project_name) /bin/sh

stop:
	docker stop $(project_name)

start:
	docker start $(project_name)
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
