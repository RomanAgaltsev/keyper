# ==============================================================================
# Code generation

.PHONY: gen-mock
gen-mock:	# Generate mocks
	go generate ./...

.PHONY: gen-proto
gen-proto:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import api/records_v1/records.proto

.PHONY: gen-buf
buf-gen:
	buf generate

# ==============================================================================
# Buf

.PHONY: buf-upd
buf-upd:
	buf dep update

.PHONY: buf-lint
buf-lint:
	buf lint

# ==============================================================================
# Local commands

.PHONY: lint
lint:
	golangci-lint run ./...
	nilaway -include-pkgs="github.com/RomanAgaltsev/keyper" ./...
	govulncheck ./...

.PHONY: lint-upd
lint-upd:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: tidy
tidy:	# Cleanup go.mod
	go mod tidy

.PHONY: fmt
fmt:
	gofumpt -l -w .
	goimports -local "github.com/RomanAgaltsev/keyper" -w .

.PHONY: test
test:	# Execute the unit tests
	go test -count=1 -race -v -timeout 30s -coverprofile cover.out ./...

.PHONY: cover
cover:	# Show the cover report
	go tool cover -html cover.out

.PHONY: update
update:	# Update dependencies as recorded in the go.mod and go.sum files
	go list -m -u all
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:	# Clean modules cache
	go clean -modcache

.PHONY: doc
doc:	# godoc
	godoc -http=:6060