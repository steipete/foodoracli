.PHONY: test fmt lint tools

test:
	go test ./...

tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

fmt:
	gofumpt -w .

lint:
	golangci-lint run
