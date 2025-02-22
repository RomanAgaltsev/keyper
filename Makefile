
.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	gofumpt -l -w .
	goimports -local "github.com/RomanAgaltsev/keyper" -w .