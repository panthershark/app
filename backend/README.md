# Getting Started

Ensure you have the latest version of these dependencies installed. The libraries will be fetched by go automatically.

- [Go](https://go.dev/doc/install): for writing code.
- [dbmate](https://github.com/amacneil/dbmate?tab=readme-ov-file#installation): database migrations
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html#)

## Run Tests

To test everything, run `gotip test ./...` or `make test` from `./backend`. There is also a helper script: `./test.sh` that sets up a container for the tests.

## Hacking

To run the backend, create a local database container, then migrate up.

1. Start a db: `podman run --name db0 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=change_me -p $(POSTGRES_TEST_PORT):5432 -d postgres:16`
2. Run either `go run cmd/server/main.go` or `make run`
3. Open [http://localhost:4011/playground](http://localhost:4011/playground)

## Sqlc

[Sqlc](https://docs.sqlc.dev/en/stable/overview/install.html#) interprets sql files and generates the database api. Queries are no longer stored in strings or runes in the code.

Install with a package manager like `pacman -S sqlc` or grab the tool from the sqlc docs, then unzip.

```bash
curl -L  -o sqlc_1.18.0_linux_amd64.tar.gz  https://github.com/kyleconroy/sqlc/releases/download/v1.18.0/sqlc_1.18.0_linux_amd64.tar.gz
tar -xf sqlc_1.18.0_linux_amd64.tar.gz .
sudo mv sqlc /usr/bin/sqlc
chmod +x /usr/bin/sqlc
```

# GraphQL Backend Server

Relies on [gqlgen](https://gqlgen.com/) which is a schema-first approach to GraphQL servers.

## Run the backend server

Run `PLAYGROUND=true gotip run ./services/server/main.go` from `./backend`. Using `PLAYGROUND=true` enables the graphql playground

# Development

The following are helpful commands for working with the project.

## Code Generation

Run `go generate ./...` from the `./backend` folder.

## Debugging with VSCode

Create a `launch.json` file and add this. This will enable debugging of your server.

```json
{
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "cmd/main.go",
      "env": {
        "PLAYGROUND": "true"
      }
    }
  ]
}
```
