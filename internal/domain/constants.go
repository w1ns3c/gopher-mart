package domain

const (
	Salt           = "testsalt"
	CookieName     = "jwt"
	UserContextKey = "user"
	InvalidUserID  = int64(-1)

	// DB tables names
	TableUsers   = "users"
	TableOrders  = "orders"
	TableBalance = "balance"
)
