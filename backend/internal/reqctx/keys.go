package reqctx

type ContextKey string

const (
	CurrentUser  ContextKey = "CurrentUser"
	AccountsApi  ContextKey = "AccountsApi"
	DatabasePool ContextKey = "DatabasePool"
)
