target=x86_64-linux-musl
export CC=zig cc -target $(target)
export CXX=zig c++ -target $(target)
export CGO_ENABLED=1
TAGS='static,osuergo,netgo'
EXTLDFLAGS="-static -Oz -s"
LDFLAGS='-linkmode=external -extldflags $(EXTLDFLAGS)'
build:
	@npm run build
	@go build -tags $(TAGS) -ldflags $(LDFLAGS) -o bin/web ./cmd/web/
	@upx bin/web
	@echo "compiled you application with all its assets to a single binary => bin/web"
