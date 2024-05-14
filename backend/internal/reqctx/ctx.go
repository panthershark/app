package reqctx

import (
	"context"

	"github.com/panthershark/app/backend/api"
	"github.com/panthershark/app/backend/db/connection"
)

type PlatformStuff struct {
	Pool   connection.Pool
	NewApi api.AccountsApiFactory
}

func GetPlatformStuff(ctx context.Context) PlatformStuff {
	return PlatformStuff{
		Pool:   ctx.Value(DatabasePool).(connection.Pool),
		NewApi: ctx.Value(AccountsApi).(api.AccountsApiFactory),
	}
}
