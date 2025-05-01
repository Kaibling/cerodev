APP_NAME=cerodev
BUILD_DIR=dist
ARCH := $(shell uname -m)
buildTime := $(shell date -u "+%Y-%m-%dT%H:%M:%S")
version := $(shell git describe --tags)

ifeq ($(ARCH),x86_64)
	ARCH := amd64
else ifeq ($(ARCH),aarch64)
	ARCH := arm64
else ifeq ($(ARCH),armv7l)
	ARCH := arm
else ifeq ($(ARCH),i686)
	ARCH := 386
else ifeq ($(ARCH),riscv64)
	ARCH := riscv64
else
$(error Unsupported architecture: $(ARCH))
endif

CURRENT_TAG := $(shell git describe --tags --abbrev=0)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Extract major, minor, patch parts
MAJOR := $(shell echo $(CURRENT_TAG) | cut -d. -f1 | sed 's/v//')
MINOR := $(shell echo $(CURRENT_TAG) | cut -d. -f2)
PATCH := $(shell echo $(CURRENT_TAG) | cut -d. -f3)

# Helper functions to increment versions
define increment_major
v$(shell echo $$(($(MAJOR) + 1))).0.0
endef

define increment_minor
v$(MAJOR).$(shell echo $$(($(MINOR) + 1))).0
endef

define increment_patch
v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH) + 1)))
endef

major:
	@git tag $(call increment_major)
	@git push origin $(call increment_major)

minor:
	@git tag $(call increment_minor)
	@git push origin $(call increment_minor)

patch:
	@git tag $(call increment_patch)
	@git push origin $(call increment_patch)

current-tag:
	@echo "Current tag: $(CURRENT_TAG)"
	@echo "Branch: $(BRANCH)"


build: .deps-ui  build-ui build-be

build-be:
	GOARCH=$(ARCH) CGO_ENABLED=0  go build -ldflags "-X main.version=${version} -X main.buildTime=${buildTime} -X main.architecture=${ARCH}" -o cerodev

run-be: build-be
	./cerodev

build-final:  build-ui 
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0  go build -ldflags "-X main.version=${version} -X main.buildTime=${buildTime} -X main.architecture=${ARCH}" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -ldflags "-X main.version=${version} -X main.buildTime=${buildTime} -X main.architecture=${ARCH}" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64

run: build
	./cerodev


.deps-ui:
	@touch .deps-ui
	apt update && apt install -y unzip 
	curl -fsSL https://bun.sh/install | bash
	export PATH=${PATH}:${HOME}/.bun/bin

build-ui:
	@base=$$(pwd); \
	rm -fr ./web/static; \
	cd ui/ && bun install && bun run build; \
	cd $$base && cp -r ./ui/dist ./web/static

run-ui: .deps-ui
	cd ui/ && bun run dev

# Lint

.deps-lint:
	@echo "Building Lint dependency..."
	@touch .deps-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.2
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

lint: .deps-lint
	gofumpt -l -w .
	govulncheck ./...
	gci write --skip-generated -s standard -s default .
	golangci-lint run

lint-ui:
	cd ui && bunx eslint src/ --format stylish

# SQLC

.deps-sql:
	@touch .deps-sql
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	
build-sql: .deps-sql
	sqlc generate

# Migration

.deps-migrate:
	@touch .deps-migrate
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate

migrate: .deps-migrate
	migrate -source file://migration/data -database "sqlite://cerodev.db" up

rollback: .deps-migrate
	migrate -source file://migration/data -database "sqlite://cerodev.db" down

# cert
cert:
	openssl req -x509 -newkey rsa:4096 -nodes \
	-keyout server.key -out server.crt -days 365 \
	-subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost"

update:
	go get -u
	go mod tidy