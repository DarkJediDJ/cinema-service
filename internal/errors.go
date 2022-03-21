package internal

import "errors"

var (
	// ErrInternalFailure ...
	ErrInternalFailure = errors.New("unable to perform your request, please try again later")

	// ErrValidationFailed ...
	ErrValidationFailed = errors.New("validation failed")
)
