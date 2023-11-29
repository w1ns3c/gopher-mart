package domain

import "time"

var (
	UserContextKey = "user"
	InvalidUserID  = int64(-1)

	// variables from .env config file
	Secret              string        //= "topsecret"
	CookieName          string        // = "jwt"
	CookieHoursLifeTime time.Duration // =  time.Hour * 4
	// DB tables names
	TableUsers   string // = "users"
	TableOrders  string // = "orders"
	TableBalance string // = "balance"
)
