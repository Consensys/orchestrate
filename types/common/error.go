package common

func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}

// Error (implement error interface)
func (err *Error) Error() string {
	return err.GetMessage()
}

func (err *Error) SetComponent(name string) *Error {
	err.Component = name
	return err
}

func (err *Error) SetCode(code uint64) *Error {
	err.Code = code
	return err
}
