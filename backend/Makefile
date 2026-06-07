.PHONY: test generate-mocks build
generate-mocks:
	go generate ./...
test:
	go test ./... -race -count=1
build:
	go build -o bin/gin-mam-server ./cmd/server
