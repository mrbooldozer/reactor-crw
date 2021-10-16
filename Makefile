BINARY_NAME=reactor-crw
BIN_DIR=bin
BUILD_PATH=${BIN_DIR}/${BINARY_NAME}
VERSION=0.0.1

.PHONY: build clean test test-coverage dep vet

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BUILD_PATH}_${VERSION}_macOS_64bit cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_PATH}_${VERSION}_Linux_64bit cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ${BUILD_PATH}_${VERSION}_Windows_32bit.exe cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BUILD_PATH}_${VERSION}_Windows_64bit.exe cmd/main.go

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
