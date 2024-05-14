package connection

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool: minimum interface for our usage. Facilitates testing abstraction.
type Pool interface {
	AcquireFunc(ctx context.Context, fn func(*pgxpool.Conn) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, fn func(pgx.Tx) error) error
}

type PgxPool struct {
	pool *pgxpool.Pool
}

func NewPgxPool(connString string) Pool {
	p := PgxPool{pool: GetDatabasePool(connString)}

	return &p
}

// AcquireFunc: returns a connection. Cleans up after execution
func (p *PgxPool) AcquireFunc(ctx context.Context, fn func(*pgxpool.Conn) error) error {
	return p.pool.AcquireFunc(ctx, fn)
}

// BeginTxFunc: executes a database transation with automatic rollback on error, automatic commit when no error
func (p *PgxPool) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, fn func(pgx.Tx) error) error {
	return p.pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		return pgx.BeginTxFunc(ctx, conn, txOptions, fn)
	})
}
