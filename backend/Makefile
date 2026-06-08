.PHONY: test generate-mocks build vendor swagger
generate-mocks:
	go generate ./...
test:
	go test ./... -race -count=1
build:
	go build -o bin/gin-mam-server ./cmd/server
vendor:
	go mod vendor
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs --parseDependency --parseInternal
