# Gofiber Template Server

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/gkampitakis/gofiber-template-server/graphs/commit-activity)
[![Ask Me Anything !](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)](https://github.com/gkampitakis/gofiber-template-server/discussions)

## Description

A template gofiber server with basic functionality and structure ready to fork and spin up a server.

## Prerequisites

For running and developing in this repository you will have to install some tools

- `npm` or `yarn` as it uses some tools like `husky` and `semantic-release`.
- After install `npm` or `yarn` you can install dependencies with executing `npm i` or `yarn` in the root
- For local development you can use `nodemon` for hot reloading. You can install it with `npm i nodemon -g` or `yarn global add nodemon`. `make dev` starts the server in hot reload.
- You will need `go`, `docker`, `golangci-lint` installed locally.

## Contents

- This repository is ready for start developing. You can find `./app/controllers` folder where route controllers are added and in `./pkg/routes/app_routes.go` you can register your new routes.

- Inside the repository you can find also `*_test.go` files containing some basic unit tests.

- In the root folder there is also a `Dockerfile` for building an image of your application.

- Also a `Makefile` is provided with some basic commands. You can type `make help` for an overview of the commands.

- `.env.example` you can create an `.env` file with your own values `cp .env.example .env`. The `.env` value is only loaded when in development.

- Swagger Documentation. If you run server and visit `localhost:8080/swagger/index.html` you can see some basic documentation for the routes.

## Usage

For starting server in hot reload mode and developing you can do it with `make dev`.

### Docker

Build Image:

```bash
docker build . -t gofiber-server
```

Run Container:

```bash
docker run --rm -p <host-port>:<server-port> --name gofiber-server gofiber-server
```

### Profiling

> For profiling you will need to have installed `graphviz`

In development profiling is enabled. You can visit `localhost:8080/debug/pprof/` and run from there the profiling. While the profiling runs you can put some stress to your server (e.g. [autocannon](https://github.com/mcollina/autocannon)). After the profiling finishes a `profile` file will be downloaded and with the command `go tool pprof -http=:8081 profile` you can visit a UI with the graphs.

### Useful commands

- Run tests `make test`
- Run linter `make lint`
- Build application `make build`
- Update swagger docs `make update-swagger`

## Licence

MIT License

Structure Now

main.go

- app/controllers
- docs/\*\*
- pkg
  - configs
  - middleware
  - routes
  - utils

Future:

- app
  - controllers/\*\*
  - middleware
  - models
- config/\*\*
- database
- routes

---

## Performance Metrics

```bash
autocannon http://localhost:8080/test/validation  -m POST -b '{"data":[]}' -d 30 -c 100

autocannon http://localhost:8080/test/validation  -m POST -d 30 -c 100 -b '{"_id":"618ecfe7ce72e1b81a44056c","index":5,"user":{"name":"gkampitakis","email":"g@mail.com","meta":{"age":10,"addresses":["gr","uk","ger"],"files":[{"type":"pdf","path":"s3","name":"CV"},{"type":"pdf","path":"s3","name":"CV"}]}},"guid":"c66e9e76-833e-4b61-9f36-0ae8407a4168","isActive":false,"balance":"$3,869.43","picture":"http://placehold.it/32x32","age":28,"eyeColor":"green","name":"Deana Blackwell","gender":"female","company":"NETPLODE","email":"deanablackwell@netplode.com","phone":"+1 (923) 512-3849","address":"682 Fayette Street, Shindler, Nebraska, 5286","about":"Enim magna proident officia quis fugiat ullamco enim. Laboris eiusmod anim aliquip ea. Sunt culpa veniam magna ullamco exercitation amet tempor mollit ullamco commodo minim consectetur Lorem voluptate. Labore cupidatat mollit enim dolore pariatur sunt pariatur ex esse anim consequat minim dolor veniam. Culpa eiusmod magna sint adipisicing non. Sunt officia veniam reprehenderit aliqua duis cillum dolore elit reprehenderit reprehenderit. Laborum pariatur do pariatur sint ullamco officia officia deserunt.\r\n","registered":"2014-07-31T03:02:15 -01:00","latitude":64.794552,"longitude":-110.078331,"tags":["proident","et","veniam","irure","quis","proident","eiusmod"],"friends":[{"id":0,"name":"Chang Rush"},{"id":1,"name":"Saundra Ellis"},{"id":2,"name":"Sherman Quinn"}],"greeting":"Hello, Deana Blackwell! You have 6 unread messages.","favoriteFruit":"banana"}'
```