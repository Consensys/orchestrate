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
	return fmt.Sprintf("%v@%v: %v", err.Hex(), err.GetComponent(), err.GetMessage())
}

// SetMessage sets error message
func (err *Error) SetMessage(format string, a ...interface{}) *Error {
	if err != nil {
		err.Message = fmt.Sprintf(format, a...)
	}
	return err
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
