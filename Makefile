export GO111MODULE=on

OUTPUT ?= dist/secret-share-web
BINARY_NAME=secret-share-web

ifdef VERSION
	LD_FLAGS="-s -w -X github.com/prog-supdex/mini-project/milestone-code/version.number=$(VERSION)"
else
	COMMIT := $(shell sh -c 'git log --pretty=format:"%h" -n 1 ')
	VERSION := $(shell sh -c 'git tag -l --sort=-version:refname "v*" | head -n1')
	LD_FLAGS="-s -w -X github.com/prog-supdex/mini-project/milestone-code/version.sha=$(COMMIT) -X github.com/prog-supdex/mini-project/milestone-code/version.number=$(VERSION)"
endif

ifndef DATA_FILE_PATH
	export DATA_FILE_PATH=./data.json
endif

# Standard build
default: build

# Install current version
install:
	go mod tidy
	go install ./...

build:
	go build -ldflags $(LD_FLAGS) -o $(OUTPUT) cmd/$(BINARY_NAME)/main.go

build-clean:
	rm -rf ./dist

# Run server
run: build
	./$(OUTPUT)

test:
	go test -count=1 -timeout=30s -race ./...

bin/golangci-lint:
	@test -x $$(go env GOPATH)/bin/golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.44.0

lint: bin/golangci-lint
	$$(go env GOPATH)/bin/golangci-lint run

fmt:
	go fmt ./...
