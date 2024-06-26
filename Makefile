.PHONY: dev-start dev-run

dev-start:
	@docker compose -f compose.dev.yml up -d

dev-run:
	@go run main.go serve
