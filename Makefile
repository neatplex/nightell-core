DIR := $(shell pwd)
GO_CI := golangci/golangci-lint:v1.59.1

.PHONY: local-up
local-up:
	@docker compose -f compose.local.yml up -d

.PHONY: local-down
local-down:
	@docker compose -f compose.local.yml down

.PHONY: local-run
local-run:
	@go run main.go serve

.PHONY: lint
lint:
	@docker run --rm -v $(DIR):/app -w /app $(GO_CI) golangci-lint run -v
