package errors

import "errors"

var (
	// gopher
	ErrGophermart = errors.New("params not initialized")

	// order add errors
	ErrAlreadyExist     = errors.New("order number already added by this user")
	ErrCreatedByAnother = errors.New("order number already added by another user")

	// order status errors
	ErrOrderNotFound = errors.New("order not found in accounting system")

	// cookie error
	ErrValueTooLong  = errors.New("cookie value too long")
	ErrInvalidCookie = errors.New("invalid cookie")

	// balance/order WithdrawErrors
	ErrNotEnoughBonuses = errors.New("user don't have so many bonuses")
	ErrWrongOrder       = errors.New("wrong ordernumber")

	// users
	ErrUserNotFoundInContext = errors.New("user not in context")
	ErrUserLogin             = errors.New("wrong password")

	ErrConfirmPassword = errors.New("password not confirmed")

	// DB
	ErrWrongResultValues = errors.New("wrong count of results")
	ErrRepoNotInit       = errors.New("repo not initialize")
	ErrDBConnect         = errors.New("can't connect to datebase")

	// handlers
	ErrMethodNotAllowed = errors.New("got not allowed HTTP method")
	ErrWrongContentType = errors.New("got not allowed content-type")
)
