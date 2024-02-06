# Step 1: Modules caching
FROM golang:1.21.6-alpine3.19 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21.6-alpine3.19 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /dist/secret-share-web cmd/secret-share-web/main.go

# Step 3: Final
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/tmp /tmp
COPY --from=builder /dist/secret-share-web /secret-share-web

EXPOSE 8080

ENTRYPOINT ["/secret-share-web"]
