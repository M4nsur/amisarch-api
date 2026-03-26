include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d amisarch-postgres

env-down:
	@docker compose down amisarch-postgres

env-cleanup:
	@read -p "Clear all environment volume files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down amisarch-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Environment files cleared"; \
	else \
		echo "Cleanup cancelled"; \
	fi

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Missing required parameter name, example: make migrate-create name=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user $(shell id -u):$(shell id -g) amisarch-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq \
		$(name)

migrate-up:
	@make migrate-action actiot=up

migrate-down:
	@make migrate-action aciot=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter action, example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm amisarch-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@amisarch-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)
