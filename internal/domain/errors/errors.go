package errors

import "errors"

var (
	// order add errors
	ErrAlreadyExist     = errors.New("order number already added by this user")
	ErrCreatedByAnother = errors.New("order number already added by another user")

	// cookie error
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")

	// balance/order WithdrawErrors
	ErrNotEnoughBonuses = errors.New("user don't have so many bonuses")
	ErrWrongOrder       = errors.New("wrong ordernumber")
)
