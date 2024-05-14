package reqctx

import (
	"context"
	"net/http"
	"os"

	"github.com/panthershark/app/backend/api"
	"github.com/panthershark/app/backend/db/connection"
)

type ReqContextOption func(ctx context.Context) context.Context

// WithDatabasePool: sets a custom DatabasePool.
func WithDatabasePool(pool connection.Pool) ReqContextOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, DatabasePool, pool)
	}
}

// WithAccountsApi: sets a custom AccountsApiFactory. The factory alloows the calling function to control the
// txn semantics of the database connection.
func WithAccountsApi(factory api.AccountsApiFactory) ReqContextOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, AccountsApi, factory)
	}
}

// defaultResolverContext: set the defaults to context options is not required.
func defaultResolverContext(ctx context.Context) context.Context {
	connString := os.Getenv("DATABASE_URL")
	pool := connection.NewPgxPool(connString)
	ctx = WithDatabasePool(pool)(ctx)
	ctx = context.WithValue(ctx, AccountsApi, api.AccountsApiFactory(api.NewAccountsApi))

	return ctx
}

// ToResolverContext: decorates a context with items needed for the GraphQL api.
func ToResolverContext(ctx context.Context, opts ...ReqContextOption) context.Context {
	ctx = defaultResolverContext(ctx)

	for _, opt := range opts {
		ctx = opt(ctx)
	}

	return ctx
}

// GraphQLMiddleware: Adds things to context that are needed for serving requests
func GraphQLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ToResolverContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
