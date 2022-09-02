export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO=$(shell which go)
export BIN=$(ROOT)/bin

export DOCKER=$(shell which docker)
export BUILD=cd $(ROOT) && $(GO) install -v -ldflags "-s -w"


default:all

.PHONY: all
all:
	$(BUILD) ./cmd/...

build:lint
	go build -o bin/main cmd/server/main.go
run:
	go run cmd/server/main.go
lint:
	golangci-lint run
compile:lint
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd cmd/server/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux cmd/server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows cmd/server/main.go

.PHONY: test
test:
	echo "***Test all of the packages***"
	go test -coverprofile=coverage.out ./...
	echo "***Calculate the coverage***"
	go tool cover -func=coverage.out


.PHONY: docker
docker: clean
	$(DOCKER) build --network host --target server -t server:${BUILD_VERSION} -f ./ci/Dockerfile .

.PHONY: clean
clean:
	@rm -f ./bin/*