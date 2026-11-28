package domain

import "errors"

// Sentinel errors represent known business failures shared across layers.
// Handlers translate these into stable HTTP status codes and JSON error codes.
var (
	ErrWorkshopNotFound = errors.New("workshop not found")
	ErrInvalidInput     = errors.New("invalid input")
)

// ValidationError carries a human-readable validation reason while still
// matching ErrInvalidInput via errors.Is().
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) Unwrap() error {
	return ErrInvalidInput
}
