package internal

import "errors"

var (
	// ErrInternalFailure creates new internal error
	ErrInternalFailure = errors.New("unable to perform your request, please try again later")

	// ErrValidationFailed creates new validation error
	ErrValidationFailed = errors.New("validation failed")

	// ErrNoSeats creates new seats error
	ErrNoSeats = errors.New("no seats")

	// ErrWrongEmail creates new email format error
	ErrWrongEmail = errors.New("wrong email format")
)
