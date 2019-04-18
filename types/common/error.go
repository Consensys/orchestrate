package common

import "fmt"

// Error (implement error interface)
func (err *Error) Error() string {
	return fmt.Sprintf("Error #%v: %v", err.Type, err.Message)
}

// Errors represent a set of error's specification
type Errors []*Error

// Error implements the error interface.
func (err Errors) Error() string {
	if len(err) == 0 {
		return ""
	}

	errors := []string{}
	for _, e := range err {
		errors = append(errors, e.Error())
	}
	return fmt.Sprintf("%v error(s): %q", len(err), errors)
}
