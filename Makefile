SHELL=/bin/bash
GOTOOLCHAIN=go1.24.1
CLIENT_NAME=keyper-cli
SERVER_NAME=keyper-srv

POSTGRES_NAME = postgres17
POSTGRES_USER = postgres
POSTGRES_PASSWORD = postgres
POSTGRES_IMAGE = postgres:17
POSTGRES_HOST = localhost:5432
POSTGRES_DB = keyper

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

# ==============================================================================
# Database

.PHONY: pg-up
pg-up: # Run Postgres docker image
	docker run --name $(POSTGRES_NAME) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -p 5432:5432 -d $(POSTGRES_IMAGE)

.PHONY: postgres-stop
postgres-stop:
	docker stop $(POSTGRES_NAME)

.PHONY: create-db
create-db:
	docker exec -it $(POSTGRES_NAME) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

.PHONY: drop-db
drop-db:
	docker exec -it $(POSTGRES_NAME) dropdb --username=$(POSTGRES_USER) $(POSTGRES_DB)

.PHONY: mig-up
mig-up:	# Apply all available migrations
	goose -dir migrations postgres "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: mig-down
mig-down: # Roll back a single migration from the current version
	goose -dir migrations postgres "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)/$(POSTGRES_DB)?sslmode=disable" down

# ==============================================================================
# Build

.PHONY: clear prep perm

build-bin: prep clients server perm

clear:
	rm -rf bin/*

prep:
	GOTOOLCHAIN=$(GOTOOLCHAIN) go mod tidy

clients:
	GOOS=linux GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -gcflags="all=-N -l" -buildvcs=false -o=bin/$(CLIENT_NAME)-linux-amd64 -o=bin/$(CLIENT_NAME) ./cmd/$(CLIENT_NAME)/...
	GOOS=windows GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -gcflags="all=-N -l" -buildvcs=false -o=bin/$(CLIENT_NAME)-windows-amd64.exe ./cmd/$(CLIENT_NAME)/...
	GOOS=darwin GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -gcflags="all=-N -l" -buildvcs=false -o=bin/$(CLIENT_NAME)-darwin-amd64 ./cmd/$(CLIENT_NAME)/...
	GOOS=darwin GOARCH=arm64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -gcflags="all=-N -l" -buildvcs=false -o=bin/$(CLIENT_NAME)-darwin-arm64 ./cmd/$(CLIENT_NAME)/...

server:
	GOOS=linux GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -gcflags="all=-N -l" -buildvcs=false -o=bin/$(SERVER_NAME)-linux-amd64 -o=bin/$(SERVER_NAME) ./cmd/$(SERVER_NAME)/...

perm:
	chmod -R +x bin

# ==============================================================================
#  Docker commands

.PHONY: dc-build
dc-build:	# Build docker compose
	docker-compose -f deployments/docker-compose.yaml -f deployments/keyper/docker-compose.keyper.yaml -f deployments/postgres/docker-compose.postgres.yaml --env-file .env build

.PHONY: dc-up
dc-up:	# Build docker compose
	docker-compose -f deployments/docker-compose.yaml -f deployments/keyper/docker-compose.keyper.yaml -f deployments/postgres/docker-compose.postgres.yaml --env-file .env up -d

.PHONY: dc-down
dc-down:	# Build docker compose
	docker-compose -f deployments/docker-compose.yaml -f deployments/keyper/docker-compose.keyper.yaml -f deployments/postgres/docker-compose.postgres.yaml --env-file .env down