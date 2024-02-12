FROM golang:1.21.6-alpine3.19 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /dist/secret-share-web cmd/secret-share-web/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /build/tmp /tmp
COPY --from=builder /dist/secret-share-web /secret-share-web

EXPOSE 8080

ENTRYPOINT ["/secret-share-web"]
