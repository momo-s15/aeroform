APP_NAME := aeroform
MODULE   := github.com/momo-s15/aeroform
VERSION  ?= dev
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE     := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
LDFLAGS  := -s -w \
	-X $(MODULE)/cmd.Version=$(VERSION) \
	-X $(MODULE)/cmd.Commit=$(COMMIT) \
	-X $(MODULE)/cmd.Date=$(DATE)

.PHONY: build install test integration e2e fmt lint clean

build:
	go build -ldflags '$(LDFLAGS)' -o $(APP_NAME) .

install:
	go install -ldflags '$(LDFLAGS)' .

test:
	go test ./...

integration:
	docker compose up -d --wait
	go test -tags integration -timeout 300s -v ./test/integration/...
	docker compose down

e2e:
	go test -tags e2e -timeout 900s -v ./test/e2e/...

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

clean:
	go clean
	rm -f $(APP_NAME)
