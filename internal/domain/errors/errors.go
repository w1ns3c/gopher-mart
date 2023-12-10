package errors

import (
	"errors"
	"fmt"
)

var (
	// gopher
	ErrGophermart = errors.New("params not initialized")

	// order add errors
	ErrOrderAlreadyExist     = errors.New("order number already added by this user")
	ErrOrderCreatedByAnother = errors.New("order number already added by another user")
	ErrOrderWrongFormat      = errors.New("order has wrong number")
	ErrOrderNotExist         = errors.New("order not exist")

	// order status errors
	ErrOrderNotFound              = errors.New("order not found in accounting system")
	ErrAccrualsNotUpdated         = errors.New("accruals info not updated")
	ErrAlreadyUpdated             = fmt.Errorf("maybe, already updated, %v", ErrAccrualsNotUpdated)
	ErrTooManyRequests            = errors.New("too many requests")
	ErrOrderNotRegisteredInRemote = errors.New("order is not registered in accruals system")

	// cookie error
	ErrValueTooLong  = errors.New("cookie value too long")
	ErrInvalidCookie = errors.New("invalid cookie")

	// balance/order WithdrawErrors
	ErrNotEnoughBonuses = errors.New("user have not enough bonuses")

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
