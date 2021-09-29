APP_NAME = gofiber-template-server
BUILD_DIR = ${PWD}/build
EXECUTABLES = docker nodemon golangci-lint
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(warning "No $(exec) in PATH")))

run: build
	./build/$(APP_NAME)

dev: __cp_env
	nodemon --exec go run main.go --signal SIGTERM

depencies:
	go mod download

clean:
	rm -rf ./build
	rm -rf .env

__cp_env:
	(cp -n .env.example .env && echo "created .env") || echo "file already exists"

lint: 
	golangci-lint run

build: clean
	cp -n .env.example .env
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

help:
	@echo 'Application: $(APP_NAME)'
	@echo
	@echo 'Usage:'
	@echo '    make update-swagger        Update swagger documentation'
	@echo '    make dependecies           Download depencies needed for running ${APP_NAME}'
	@echo '    make dev                   Run project with hot reload (needs nodemon installed)'
	@echo '    make run                   Build and run project for development (enabled pprof)'
	@echo '    make test                  Run unit tests (you can pass -v for logs or debugging)'
	@echo '    make clean                 Delete ./build and .env'
	@echo '    make build                 Build application and copy the .env file'
	@echo

test: lint
	go test ./... -count=1 -race

update-swagger:
	~/go/bin/swag init ./...