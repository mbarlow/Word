.PHONY: help build dev swagger test clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	go build -o bin/word-service ./cmd/server

dev: ## Start dev server (air + browser-sync)
	./scripts/dev.sh

swagger: ## Generate Swagger docs to swagger/
	swag init -g cmd/server/main.go --output swagger --parseInternal --parseDependency

test: ## Run tests
	go test ./... -v -race -timeout 120s

clean: ## Remove build artifacts
	rm -rf bin tmp swagger
