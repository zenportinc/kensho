# Helpers

.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: mod
mod: ## Download GO dependencies.
	go mod vendor

# Builds

.PHONY: build-go
build-go: mod ## Build the Go binary
	GOPRIVATE=$(GO_PRIVATE) CGO_ENABLED=0 go build .

# Application

# Tools

.PHONY: format
format:
	go run github.com/daixiang0/gci@latest write --skip-generated --skip-vendor .

.PHONY: linter
linter: mod ## Run golangci-lint
	docker compose up --force-recreate --exit-code-from linter linter

# Tests

.PHONY: unit-tests
unit-tests: mod ## Run unit tests
	go test ./...	\
		-shuffle=on 			\
		-tags=unit,benchmark	\
		-covermode=atomic 		\
		-bench=. 				\
		-count=1				\
		-race 					\
        -failfast 				\
      	-coverprofile=coverage.out && \
    go tool cover -html=coverage.out -o=coverage.html
