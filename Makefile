PROJECT_NAME := "book"
MAIN_FILE_PATH := "main.go"
MOD_DIR := $(shell go env GOPATH)
PKG := "github.com/ducketlab/$(PROJECT_NAME)"
IMAGE_PREFIX := "github.com/ducketlab/book"

all: build

dep: ## Get the dependencies
	@go mod download

install: ## Install dependence go package
	@go install github.com/golang/protobuf/protoc-gen-go
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install github.com/ducketlab/mingo/cmd/protoc-gen-go-http

codegen: # Init Service
	@protoc -I=. -I${MOD_DIR}/src --go-ext_out=. --go-ext_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} --go-http_out=. --go-http_opt=module=${PKG} pkg/*/pb/*.proto
	@go generate ./...

build: dep ## Build the binary file
	@go fmt ./...
	@bash ./scripts/build.sh local dist/${PROJECT_NAME} ${MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

linux: ## Linux build
	@bash ./scripts/build.sh linux dist/${PROJECT_NAME} ${MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

run: install codegen dep build ## Run Server
	@./dist/${PROJECT_NAME} start

clean: ## Remove previous build
	@go clean .
	@rm -f dist/${PROJECT_NAME}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
