build:
	go build -o ./bin/calendar ./cmd/calendar

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
lint: install-lint-deps
	golangci-lint run ./...

test-db:
	./scripts/db-dev.sh
.PHONY: build test lint
