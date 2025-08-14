include .env

LDFLAGS := -X main.buildTime=$(shell date +%Y-%m-%d,%H:%M:%S) \
			-X main.buildVersion=${WANWU_WORKFLOW_VERSION} \
			-X main.gitCommitID=$(shell git --git-dir=../.git rev-parse HEAD) \
			-X main.gitBranch=$(shell git --git-dir=../.git for-each-ref --format='%(refname:short)->%(upstream:short)' $(shell git --git-dir=../.git symbolic-ref -q HEAD)) \
			-X main.builder=$(shell git config user.name)

build-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/workflow-wanwu

build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/workflow-wanwu

check:
	go vet ./...
	go fmt ./...
	docker run --rm -t -v $(PWD):/app -w /app golangci/golangci-lint:v1.64.8 bash -c 'golangci-lint run -v --timeout 3m'

init:
	cp ../docker/volumes/mysql/schema.sql configs/
	go mod tidy
	go mod vendor