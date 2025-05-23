.PHONY: help vendor generate test

help: ## Show help
	@echo "\n\033[1mAvailable commands:\033[0m\n"
	@@awk 'BEGIN {FS = ":.*##";} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

vendor: ## Install dependencies using vendoring
	go mod vendor

migrate-up: ## Apply all available migrations
	go run ./cmd/migrations -up

migrate-down: ## Rollback the most recent migration
	go run ./cmd/migrations -down

migrate-version: ## Show the current migration version
	go run ./cmd/migrations -version

generate: ## Generate
	go generate ./...

test: generate ## Testing
	go test ./...