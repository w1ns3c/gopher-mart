package errors

import "errors"

var (
	// order add errors
	ErrAlreadyExist     = errors.New("order number already added by this user")
	ErrCreatedByAnother = errors.New("order number already added by another user")

	// order status errors
	ErrOrderNotFound = errors.New("order not found in accounting system")

	// cookie error
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")

	// balance/order WithdrawErrors
	ErrNotEnoughBonuses = errors.New("user don't have so many bonuses")
	ErrWrongOrder       = errors.New("wrong ordernumber")

	// users
	ErrUserNotFoundInContext = errors.New("user not in context")
	ErrUserLogin             = errors.New("wrong password")

	ErrConfirmPassword = errors.New("password not confirmed")

	// DB
	ErrWrongResultValues = errors.New("wrong count of results")

	// handlers
	ErrMethodNotAllowed = errors.New("got not allowed HTTP method")
	ErrWrongContentType = errors.New("got not allowed content-type")
)
