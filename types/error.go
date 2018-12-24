package types

const (
	// ErrorTypeLoad is used when protobuffer loading fails
	ErrorTypeLoad = 1

	// ErrorTypeUnknown is used when error is unknown
	ErrorTypeUnknown = 2

	// ErrorTypeDone is used when context Timeout or is Cancelled
	ErrorTypeDone = 2
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
