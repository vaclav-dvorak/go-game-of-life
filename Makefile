GIT_REV?=$$(git rev-parse --short HEAD)
DATE?=$$(date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=$$(git describe --tags --always)
LDFLAGS="-s -w -X main.version=$(VERSION)-$(GIT_REV) -X main.date=$(DATE)"

default: help

cover:  ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build:  ## Builds the game
	@CGO_ENABLED=1 go build -ldflags=$(LDFLAGS) -o ./bin/gol main.go

buildwin: ## Builds game for Win64
	@CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags=$(LDFLAGS) -o ./bin/gol.exe main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
