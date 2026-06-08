.PHONY: test generate-mocks build vendor
generate-mocks:
	go generate ./...
test:
	go test ./... -race -count=1
build:
	go build -o bin/gin-mam-server ./cmd/server
vendor:
	go mod vendor
