.PHONY: up down setup clean

up:
	docker-compose -f infrastructure/docker-compose.yml up -d

down:
	docker-compose -f infrastructure/docker-compose.yml down

setup:
	@echo "Creating .env files if they don't exist..."
	@touch accio-api/.env
	@touch accio-mcp/.env
	@echo "Setup complete."

clean:
	docker-compose -f infrastructure/docker-compose.yml down -v