.FORCE:
.PHONY: build test-only test docker-build watch .FORCE

BIN_NAME = api

COMMIT     = $(shell git rev-parse HEAD)
BRANCH     = $(shell git branch --show-current | awk '{print substr(tolower($$0),0,40)}')
VERSION    = $(shell git describe --tags --match=v* 2>/dev/null || git rev-parse HEAD)
BUILD_DATE = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COVERPKG   = ./cmd/...,./internal/...,./migrations/...

all: lint test tools build

tidy:
	go mod tidy

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/TimRutte/api/internal/version.Version=${VERSION} -X github.com/TimRutte/api/internal/version.GitCommit=${COMMIT} -X github.com/TimRutte/api/internal/version.BuildDate=${BUILD_DATE}" -o ${BIN_NAME} ./cmd/api

build-darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "-X github.com/TimRutte/api/internal/version.Version=${VERSION} -X github.com/TimRutte/api/internal/version.GitCommit=${COMMIT} -X github.com/TimRutte/api/internal/version.BuildDate=${BUILD_DATE}" -o ${BIN_NAME} ./cmd/api

build-windows:
	CGO_ENABLED=0 GOOS=windows go build -ldflags "-X github.com/TimRutte/api/internal/version.Version=${VERSION} -X github.com/TimRutte/api/internal/version.GitCommit=${COMMIT} -X github.com/TimRutte/api/internal/version.BuildDate=${BUILD_DATE}" -o ${BIN_NAME}.exe ./cmd/api

run:
	./${BIN_NAME}

test-only:
	go test -coverpkg ${COVERPKG} -coverprofile coverage.out ./...

test: test-only
	govulncheck ./...

watch:
	~/go/bin/gow test ./...

lint:
	golangci-lint run --fix

coverage: test-only coverage-html

coverage-html:
	go tool cover -html=coverage.out

coverage-func:
	go tool cover -func=coverage.out

clean:
	@if [ -f "${BIN_NAME}" ]; then rm "${BIN_NAME}"; fi

githooks-install:
	pre-commit install

githooks-uninstall:
	pre-commit uninstall

tools: .FORCE
	@echo "getting tools"
	go generate tools/*.go

docker-compose-up:
	docker compose up -d

docker-build: tools
	docker build -f docker/Dockerfile-app -t "${BIN_NAME}:latest" .

docker-run:
	docker run -p 50052:50052 -p 8099:8099 --name "${BIN_NAME}_container" "${BIN_NAME}:latest"

docker-stop:
	docker stop "${BIN_NAME}_container" || true
	docker rm "${BIN_NAME}_container" || true
