package graph_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/panthershark/app/backend/api"
	"github.com/panthershark/app/backend/db/dbc"
	"github.com/panthershark/app/backend/internal/graph"
	"github.com/panthershark/app/backend/internal/reqctx"
	"github.com/panthershark/app/backend/internal/testing/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func withMockApi(a api.AccountsApi) reqctx.ReqContextOption {
	factory := func(querier dbc.Querier, ctx *context.Context) api.AccountsApi {
		return a
	}

	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, reqctx.AccountsApi, api.AccountsApiFactory(factory))
	}
}

// withAquireFuncPool: adds a mock pool that expects to call pool.AcquirFunc()
// which does not open a database txn
// This is a convenience for common cases. If you need something special, use reqctx.WithDatabasePool()
func withAquireFuncPool(ctrl *gomock.Controller) reqctx.ReqContextOption {
	pool := mocks.NewMockPool(ctrl)
	pool.EXPECT().AcquireFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(*pgxpool.Conn) error) error {
		return fn(nil)
	})

	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, reqctx.DatabasePool, pool)
	}
}

// withBeginTxPool: adds a mock pool that expects to call pool.BeginTx()
// which opens a database txn.
// This is a convenience for common cases. If you need something special, use reqctx.WithDatabasePool()
func withBeginTxPool(ctrl *gomock.Controller) reqctx.ReqContextOption {
	pool := mocks.NewMockPool(ctrl)
	pool.EXPECT().BeginTxFunc(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, txOptions pgx.TxOptions, fn func(pgx.Tx) error) error {
		return fn(nil)
	})
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, reqctx.DatabasePool, pool)
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accountsApi := mocks.NewMockAccountsApi(ctrl)
	accountsApi.EXPECT().GetUserByEmail(gomock.Eq("hulk@hogan.com")).Return(
		&api.Account{
			ID:    uuid.MustParse("8b8c0b77-837f-400b-a943-093f92096580"),
			Email: "hulk@hogan.com",
			Person: &api.Person{
				ID:        uuid.MustParse("00839969-2b43-4652-aa84-1effcfd477d6"),
				AccountID: uuid.MustParse("8b8c0b77-837f-400b-a943-093f92096580"),
				FirstName: "Hulk",
				LastName:  "Hogan",
			},
		}, nil)

	ctx := reqctx.ToResolverContext(
		context.Background(),
		withAquireFuncPool(ctrl),
		withMockApi(accountsApi),
	)
	resolver := graph.Resolver{}

	got, err := resolver.Query().GetUser(ctx, "hulk@hogan.com")
	assert.NoError(t, err)
	// cupaloy.SnapshotT(t, redactif.Redact(&got, "snapshot"))
	cupaloy.SnapshotT(t, got)
}

func TestCreateUser(t *testing.T) {
	input := api.UserCreateInput{
		Email:     "hulk@hogan.com",
		FirstName: "Hulk",
		LastName:  "Hogan",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	accountsApi := mocks.NewMockAccountsApi(ctrl)
	accountsApi.EXPECT().CreateUser(gomock.Eq(input)).Return(uuid.MustParse("5732374c-479c-4d36-8df6-26d650e44aa2"), nil)

	ctx := reqctx.ToResolverContext(
		context.Background(),
		withBeginTxPool(ctrl),
		withMockApi(accountsApi),
	)

	resolver := graph.Resolver{}
	got, err := resolver.Mutation().CreateUser(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, "5732374c-479c-4d36-8df6-26d650e44aa2", got.String())
}
