BINARY_NAME=reactor-crw
BIN_DIR=bin
BUILD_PATH=${BIN_DIR}/${BINARY_NAME}
VERSION=0.0.1

.PHONY: build clean test test-coverage dep vet

build:
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_PATH}_${VERSION}_macOS_64bit cmd/main.go && upx --best --lzma ${BUILD_PATH}_${VERSION}_macOS_64bit
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_PATH}_${VERSION}_Linux_64bit cmd/main.go && upx --best --lzma ${BUILD_PATH}_${VERSION}_Linux_64bit
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_PATH}_${VERSION}_Windows_64bit.exe cmd/main.go && upx --best --lzma ${BUILD_PATH}_${VERSION}_Windows_64bit.exe

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
