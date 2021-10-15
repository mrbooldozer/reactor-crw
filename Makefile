BINARY_NAME=reactor-crw
BIN_DIR=bin
BUILD_PATH=${BIN_DIR}/${BINARY_NAME}

.PHONY: build clean test test-coverage dep vet

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BUILD_PATH}-darwin cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_PATH} cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BUILD_PATH}.exe cmd/main.go

clean:
	@go clean
	@rm -rf ${BIN_DIR}

test:
	go test ./... --tags unit

test-coverage:
	go test ./... --tags unit -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet
