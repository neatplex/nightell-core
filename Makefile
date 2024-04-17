.PHONY: dev-start dev-run

dev-start:
	@docker compose -f docker-compose.db.yml up -d

dev-run:
	@go run main.go serve
