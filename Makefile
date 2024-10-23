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
	--build.include_ext "go,tmpl" \
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
	@go build -o bin/snippets ./cmd/web/
	@echo "compiled you application with all its assets to a single binary => bin/snippets"
