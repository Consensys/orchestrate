package error

import (
	"fmt"
)

// New returns an error with given message
func New(msg string) *Error {
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
	if err != nil {
		err.Component = name
	}
	return err
}

// ExtendComponent extend the component
func (err *Error) ExtendComponent(name string) *Error {
	if err.GetComponent() == "" {
		_ = err.SetComponent(name)
	} else {
		err.Component = fmt.Sprintf("%v.%v", name, err.Component)
	}
	return err
}

// Hex returns error code in HEX reprensatation
func (err *Error) Hex() string {
	return fmt.Sprintf("%05X", err.GetCode())
}

// SetCode sets error code
func (err *Error) SetCode(code uint64) *Error {
	if err != nil {
		err.Code = code
	}
	return err
}

// Errorf creates an error according to a format specifier
func Errorf(format string, a ...interface{}) *Error {
	return New(fmt.Sprintf(format, a...))
}

// FromError cast a golang error into an Error pointer
//
// if `err` is an internal error then it is returned
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		return New(err.Error())
	}

	return e
}
