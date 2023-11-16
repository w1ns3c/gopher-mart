package localerrors

import "errors"

var (
	ErrConfirmPassword = errors.New("password not confirmed")
)
