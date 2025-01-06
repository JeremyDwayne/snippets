build:
	@npm run build
	@go build -o bin/app_prod ./cmd/web/
	@upx bin/app_prod
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"
