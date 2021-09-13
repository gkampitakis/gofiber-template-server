APP_NAME = gofiber-template-server
BUILD_DIR = ${PWD}/build
EXECUTABLES = docker nodemon
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(warning "No $(exec) in PATH")))

run: build
	./build/$(APP_NAME)

dev:
	(cp -n .env.example .env && echo "created .env") || echo "file already exists"
	nodemon --exec go run main.go --signal SIGTERM

clean:
	rm -rf ./build
	rm -rf .env

build: clean 
	cp -u .env.example .env
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

help:
	@echo 'Application: $(APP_NAME)'
	@echo
	@echo 'Usage:'
	@echo '    make dev										Run project with hot reload ( needs nodemon installed )
	@echo '    make run                   Running project for development (enabled pprof)'
	@echo '    make clean                 Delete ./build and .env'
	@echo '    make build                 Build application and copy the .env file'
	@echo '    make unit-test             Run unit tests'
	@echo '    make integration-tests     Run integration tests (First you need to start docker make docker-start)'
	@echo

unit-tests:
	go test -run Unit ./...

docker-start:
	@echo 'After running tests you can stop docker with `make docker-stop`'
	docker-compose up -d

docker-stop:
	docker-compose down

integration-tests:
	go test -run Integration -p=1 -v