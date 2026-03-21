include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d amisarch-postgres

env-down:
	docker compose down amisarch-postgres

env-cleanup:
	@read -p "Clear all environment volume files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down amisarch-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Environment files cleared"; \
	else \
		echo "Cleanup cancelled"; \
	fi
