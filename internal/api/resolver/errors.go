package resolver

import "errors"

// API errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrNotLogged       = errors.New("user is not logged")
	ErrMainEqualsCrush = errors.New("mainID and crushID must not be the same")
	ErrFetchFailed     = errors.New("failed to fetch the documents")
)
