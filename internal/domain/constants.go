package domain

import "time"

var (
	UserContextKey = "user"
	InvalidUserID  = int64(-1)

	// variables from .env config file
	Secret              string        = "topsecret"
	LogLevel            string        = "DEBUG"
	CookieName          string        = "jwt"
	CookieHoursLifeTime time.Duration = time.Hour * 4

	WorkersCount  uint          = 4
	RetryTimer    time.Duration = time.Second * 10
	RetryAttempts uint          = 2

	// DB tables names
	TableUsers     string = "users"
	TableOrders    string = "orders"
	TableBalance   string = "balance"
	TableWithdraws string = "withdraws"
)
