.PHONY: test build tidy fmt vet

test: tidy fmt vet

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -o url .
