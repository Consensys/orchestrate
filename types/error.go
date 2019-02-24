package types

import "fmt"

const (
	// ErrorTypeLoad is used when protobuffer loading fails
	ErrorTypeLoad = 1

	// ErrorTypeUnknown is used when error is unknown
	ErrorTypeUnknown = 2

	// ErrorTypeDone is used when context Timeout or is Cancelled
	ErrorTypeDone = 2

	// ErrorTypeNonce when error occurs on nonce handler
	ErrorTypeNonce = 8
)

// Error represents a error's specification.
type Error struct {
	Err  error
	Type uint64
}

// Error implements the error interface.
func (msg Error) Error() string {
	return msg.Err.Error()
}

// Errors represent a set of error's specification.
type Errors []*Error

// Error implements the error interface.
func (err Errors) Error() string {
	return fmt.Sprint(err)
}
