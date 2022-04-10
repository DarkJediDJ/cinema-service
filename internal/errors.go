package internal

import "errors"

var (
	// ErrInternalFailure creates new internal error
	ErrInternalFailure = errors.New("unable to perform your request, please try again later")

	// ErrValidationFailed creates nev validation error
	ErrValidationFailed = errors.New("validation failed")

	// ErrNoSeats creates nev seats error
	ErrNoSeats = errors.New("no seats")
)
