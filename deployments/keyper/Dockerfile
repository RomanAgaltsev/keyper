FROM golang:1.24.1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /keper-srv ./cmd/keper-srv/main.go

FROM alpine
WORKDIR /
COPY --from=0 /keper-srv /keper-srv
EXPOSE 8080
ENTRYPOINT ["/keper-srv"]