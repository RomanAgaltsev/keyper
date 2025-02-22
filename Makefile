.PHONY: gen-mock
gen-mock:	# Generate mocks
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run ./...
	nilaway -include-pkgs="github.com/RomanAgaltsev/keyper" ./...

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

.PHONY: doc
doc:	# godoc
	godoc -http=:6060