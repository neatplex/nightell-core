.PHONY: dev-start dev-stop dev-run

dev-start:
	@docker compose -f compose.dev.yml up -d

dev-stop:
	@docker compose -f compose.dev.yml down

dev-run:
	@go run main.go serve
