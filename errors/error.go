package errors

import (
	"fmt"

	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Errorf creates an error according to a format specifier
func Errorf(format string, a ...interface{}) *ierror.Error {
	err := InternalError(fmt.Sprintf(format, a...))
	return err
}

// FromError cast a golang error into an internal Error
//
// if `err` is an internal then it is returned
func FromError(err error) *ierror.Error {
	if err == nil {
		return nil
	}

	ierr, ok := err.(*ierror.Error)
	if !ok {
		return Errorf(err.Error())
	}

	return ierr
}

// isErrorClass returns whether code belongs to a base error class
func isErrorClass(code, base uint64) bool {
	return (base^code)&(255<<12+15<<8&base) == 0
}
