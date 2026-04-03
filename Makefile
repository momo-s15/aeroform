APP_NAME := aeroform
GO_FILES := ./...

.PHONY: build test fmt lint clean

build:
	go build ./...

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	@echo "lint not configured yet"

clean:
	go clean
