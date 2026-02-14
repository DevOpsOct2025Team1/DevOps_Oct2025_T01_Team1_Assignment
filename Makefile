.PHONY: help build test coverage serve docker-build docker-up docker-down clean

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build:
	./nx run-many --target=build --all

test:
	./nx run-many --target=test --all

coverage:
	./nx run-many --target=coverage --all

serve:
	./nx run-many --target=serve --all --parallel=3

docker-build:
	./nx run-many --target=docker-build --all

docker-up:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

clean:
	rm -rf dist/
	rm -rf logs/

dev-deps:
	@echo "Checking dependencies..."
	@which go > /dev/null || (echo "Go is not installed" && exit 1)
	@which docker > /dev/null || (echo "Docker is not installed" && exit 1)
	@which kubectl > /dev/null || (echo "kubectl is not installed" && exit 1)
	@echo "All dependencies are installed!"

setup-env:
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file"; fi
	@if [ ! -f services/user-service/.env ]; then cp apps/user-service/.env.example apps/user-service/.env; echo "Created services/user-service/.env"; fi
	@if [ ! -f services/auth-service/.env ]; then cp apps/auth-service/.env.example apps/auth-service/.env; echo "Created services/auth-service/.env"; fi
	@echo "Environment files created. Please update them with your actual credentials."