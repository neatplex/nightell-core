DIR := $(shell pwd)
GO_CI := golangci/golangci-lint:v1.59.1

.PHONY: dev-start
dev-start:
	@docker compose -f compose.dev.yml up -d

.PHONY: dev-stop
dev-stop:
	@docker compose -f compose.dev.yml down

.PHONY: dev-run
dev-run:
	@go run main.go serve

.PHONY: lint
lint:
	@docker run --rm -v $(DIR):/app -w /app $(GO_CI) golangci-lint run -v
