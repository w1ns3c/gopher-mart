package errors

import "errors"

var (
	// order add errors
	ErrAlreadyExist     = errors.New("order number already added by this user")
	ErrCreatedByAnother = errors.New("order number already added by another user")
)
