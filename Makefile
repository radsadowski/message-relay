SHELL := /bin/bash
.PHONY: help

APP_NAME = outbox 

start: ## start the outbox service
	@docker compose -f docker/docker-compose.yml up --build --wait

stop: ## stop the outbox service
	@docker compose -f docker/docker-compose.yml down -t 1

help: ## show help
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

build: ## build outbox binary
	@go build -o dist/$(APP_NAME)

install: build ## installs outbox binary - only for development purposes
	@sudo cp dist/$(APP_NAME) /usr/local/bin/

tests: ## run unit tests
	@go test -coverprofile=/tmp/outbox.out ./pkg/...

coverage: tests ## run tests with coverage
	@go tool cover -html=/tmp/outbox.out

static-analysis:  ## run static analysis
	@echo "Static analysis..."
	@golangci-lint run --config .golangci.yml --out-format=colored-line-number --concurrency 8

fixer: get-golangci-config ## run static analysis
	@echo "Fixer..."
	@golangci-lint run --config .golangci.yml --out-format=colored-line-number --concurrency 8 --fix
