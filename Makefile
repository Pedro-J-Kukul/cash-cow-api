# Environment File
include .envrc

# Migration Commands
.PHONY : migrate/new migrate/up migrate/down migrate/reset migrate/fix migrate/down-1 migrate/up-1
migrate/new:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'make migrate/new name=your_migration_name'"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate/up:
	@echo "Applying migrations to database at $(DB_DSN)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" up

migrate/down:
	@echo "Reverting migrations on database at $(DB_DSN)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" down

migrate/reset:
	@echo "Resetting database at $(DB_DSN)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" drop -f
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" up

migrate/fix:
	@echo "Checking migration status..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" version > /tmp/migrate_version 2>&1 || echo "No version found, assuming fresh database."
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" force $$version; \
		echo "Running Down migration to revert dirty state..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" down 1; \
		echo "Dirty migration fixed."; \
		echo "Running Up migration to apply changes..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" up 1; \
		echo "Migration reapplied successfully."; \
	else \
		echo "Database is clean. No action needed."; \
	fi
	@rm -f /tmp/migrate_version

migrate/down-1:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" down 1

migrate/up-1:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN)" up 1

# Seed Commands
.PHONY : seed/new seed/up seed/down seed/reset
seed/new:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a seed name using 'make seed/new name=your_seed_name'"; \
		exit 1; \
	fi
	@echo "Creating new seed: $(name)"
	migrate create -ext sql -dir $(SEEDS_PATH) -seq $(name)_seed

seed/up:
	@echo "Applying seeds to database at $(DB_DSN)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN)" up

seed/down:
	@echo "Reverting seeds on database at $(DB_DSN)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN)" down

seed/reset:
	@echo "Resetting seeds on database at $(DB_DSN)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN)" drop -f
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN)" up

# Migration Commands for Test Database
.PHONY : migrate/test-new migrate/test-up migrate/test-down migrate/test-reset migrate/test-fix
migrate/test-new:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'make migrate/test-new name=your_migration_name'"; \
		exit 1; \
	fi
	@echo "Creating new test migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
migrate/test-up:
	@echo "Applying migrations to test database at $(DB_DSN_TEST)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" up
migrate/test-down:
	@echo "Reverting migrations on test database at $(DB_DSN_TEST)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" down
migrate/test-reset:
	@echo "Resetting test database at $(DB_DSN_TEST)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" drop -f
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" up
migrate/test-fix:
	@echo "Checking test database migration status..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" version > /tmp/migrate_test_version 2>&1 || echo "No version found, assuming fresh database."
	@cat /tmp/migrate_test_version
	@if grep -q "dirty" /tmp/migrate_test_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_test_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" force $$version; \
		echo "Running Down migration to revert dirty state..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" down 1; \
		echo "Dirty migration fixed."; \
		echo "Running Up migration to apply changes..."; \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_DSN_TEST)" up 1; \
		echo "Migration reapplied successfully."; \
	else \
		echo "Test database is clean. No action needed."; \
	fi
	@rm -f /tmp/migrate_test_version


# Seeding Commands for Test Database
.PHONY : seed/test-new seed/test-up seed/test-down seed/test-reset
seed/test-new:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a seed name using 'make seed/test-new name=your_seed_name'"; \
		exit 1; \
	fi
	@echo "Creating new test seed: $(name)"
	migrate create -ext sql -dir $(SEEDS_PATH) -seq $(name)_seed
seed/test-up:
	@echo "Applying seeds to test database at $(DB_DSN_TEST)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN_TEST)" up
seed/test-down:
	@echo "Reverting seeds on test database at $(DB_DSN_TEST)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN_TEST)" down
seed/test-reset:
	@echo "Resetting seeds on test database at $(DB_DSN_TEST)"
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN_TEST)" drop -f
	migrate -path $(SEEDS_PATH) -database "$(DB_DSN_TEST)" up



# Database Commansd
.PHONY: db/login db/test db/postgres db/dump 
db/login:
	@echo "Connecting to database at $(DB_DSN)"
	@psql $(DB_DSN)
db/test:
	@echo "Connecting to test database at $(DB_DSN_TEST)"
	@psql $(DB_DSN_TEST)
db/postgres:
	@echo "Connecting to PostgreSQL as user 'postgres'"
	@sudo -u postgres psql
db/dump:
	@echo "Dumping database at $(DB_NAME) to $(DB_DUMP_PATH)"
	@pg_dump -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -s -F p -E UTF-8 -f $(DB_DUMP_PATH)
	@echo "Database dump completed."
# Helpers
.PHONY: help/migrations help/seeds help/setup-db-migrations

help/migrations:
	@if [ ! -d "$(MIGRATIONS_PATH)" ]; then \
		echo "Creating migrations directory at $(MIGRATIONS_PATH)"; \
		mkdir -p $(MIGRATIONS_PATH); \
	else \
		echo "Migrations directory already exists at $(MIGRATIONS_PATH)"; \
	fi

help/seeds:
	@if [ ! -d "$(SEEDS_PATH)" ]; then \
		echo "Creating seeds directory at $(SEEDS_PATH)"; \
		mkdir -p $(SEEDS_PATH); \
	else \
		echo "Seeds directory already exists at $(SEEDS_PATH)"; \
	fi	

help/setup-db-migrations: 
	@echo "Creating Users Table"
	@make migrate/new name=create_users_table
	@make migrate/new name=create_tokens_table
	@make migrate/new name=create_permissions_table
	@make migrate/new name=create_users_permissions_table
	@echo "Creating Cattle Table"
	@make migrate/new name=create_breeds_table
	@make migrate/new name=create_cattle_table
	@echo "Creating Locations Table"
	@make migrate/new name=create_districts_table
	@make migrate/new name=create_areas_table
	@echo "Creating Listings Table"
	@make migrate/new name=create_listings_table
	@make migrate/new name=create_listings_catle_table
	@make migrate/new name=listing_prices_table
	@echo "Database migration setup complete."