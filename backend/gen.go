package pkg

//go:generate sqlc generate
//go:generate go run go.uber.org/mock/mockgen -source=./db/dbc/querier.go -destination=./internal/testing/mocks/querier.go -package=mocks
//go:generate go run go.uber.org/mock/mockgen -source=./db/connection/pool.go -destination=./internal/testing/mocks/pool.go -package=mocks
//go:generate go run go.uber.org/mock/mockgen -source=./api/accounts.go -destination=./internal/testing/mocks/accounts.go -package=mocks
