package connection

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pools map[string]*pgxpool.Pool = map[string]*pgxpool.Pool{}

var poolLock = sync.Mutex{}

// GetDatabasePool: returns a database connection pool
func GetDatabasePool(connectionString string) *pgxpool.Pool {
	poolLock.Lock()
	defer poolLock.Unlock()

	pool, ok := pools[connectionString]
	if !ok {
		p, err := pgxpool.New(context.Background(), connectionString)
		if err != nil {
			log.Panic(err)
		}

		pools[connectionString] = p
		pool = p
	}

	return pool
}

// CloseDatabasePools: closes all open pools.
func CloseDatabasePools() {
	poolLock.Lock()
	defer poolLock.Unlock()

	for k, pool := range pools {
		pool.Close()
		delete(pools, k)
	}
}
