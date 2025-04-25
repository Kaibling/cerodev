bun := /root/.bun/bin/bun
APP_NAME=cerodev
BUILD_DIR=dist


build: ui-deps  build-ui build-be

build-be:
	GOARCH=arm64 CGO_ENABLED=0  go build -o cerodev


build-final: ui-deps  build-ui 
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0  go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64

run: build
	./cerodev

ui-deps:
	apt update && apt install -y unzip 
	curl -fsSL https://bun.sh/install | bash


lint-deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.2
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

run-ui: ui-deps
	cd ui/ && $(bun) run dev


lint:
	gofumpt -l -w .
	govulncheck ./...
	gci write --skip-generated -s standard -s default .
	golangci-lint run

build-ui:
	@base=$$(pwd); \
	rm -fr ./web/static; \
	cd ui/ && 	$(bun) install && $(bun) run build; \
	cd $$base && cp -r ./ui/dist ./web/static

deps:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate
	
build-sql:
	sqlc generate

migrate:
	migrate -source file://migration/data -database "sqlite://cerodev.db" up

rollback:
	migrate -source file://migration/data -database "sqlite://cerodev.db" down
