package errors

import (
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Encoding errors (code 01XXX)
//
// Encoding errors are raised when an error is encountered while
var encodingErrCode = []byte{0x10, 0x00}

// EncodingError creates an encoding error
func EncodingError(e error) *err.Error {
	return err.FromError(e).SetCode(encodingErrCode)
}
