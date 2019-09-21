package resolver

import "errors"

// API errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrNotLogged       = errors.New("user is not logged")
	ErrMismatchMainID  = errors.New("mainID must be the profile ID of the user")
	ErrMainEqualsCrush = errors.New("mainID and crushID must not be the same")
)
