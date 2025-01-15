# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

MAIN_PATH = tmp/main

server:
	@go run github.com/air-verse/air@latest \
	--build.cmd "go build --tags dev -o ${MAIN_PATH} ./cmd/web/" --build.bin "${MAIN_PATH}" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,tmpl,css,js" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--screen.clear_on_rebuild true \
	--log.main_only true

watch-assets:
	@npx tailwindcss -i ui/static/css/custom.css -o ui/static/css/style.css --watch

watch-esbuild:
	@npx esbuild ui/static/js/custom.js --bundle --outfile=ui/static/js/index.js --watch

dev:
	@make -j5 server watch-assets watch-esbuild

build:
	@npx tailwindcss -i ui/static/css/custom.css -o ui/static/css/style.css
	@npx esbuild ui/static/js/custom.js --bundle --outfile=ui/static/js/index.js
	@go build -o bin/app_prod cmd/web/main.go
	@upx bin/app_prod
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"

db-status:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) status

db-reset:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) reset

db-down:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) down

db-up:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) up

db-migration-create:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_NAME) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) create $(filter-out $@,$(MAKECMDGOALS)) sql
