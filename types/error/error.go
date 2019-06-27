package error

// NewError creates a new error
func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}

// Error (implement error interface)
func (err *Error) Error() string {
	return err.GetMessage()
}

// SetComponent set component
func (err *Error) SetComponent(name string) *Error {
	err.Component = name
	return err
}

// SetCode sets error code
func (err *Error) SetCode(code []byte) *Error {
	err.Code = code
	return err
}

// FromError cast an error tinto an Error
//
// if `err` is already an internal error then it is returned
func FromError(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		return NewError(err.Error())
	}
	return e
}
