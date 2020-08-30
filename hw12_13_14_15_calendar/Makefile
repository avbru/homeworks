build:
	go build -o ./bin/calendar ./cmd/calendar

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.30.0

lint: 
	gofumpt -s -w .
	golangci-lint run ./...

test-db:
	./scripts/db-dev.sh
.PHONY: build test lint
