package errors

import (
	"fmt"

	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Configuration errors (code 01XXX)
//
// Configuration errors are raised when an error is encountered validating a configuration format

var configErrCode = []byte{0x10, 0x00}

// ConfigError creates a a configuration error
func ConfigError(e error) *err.Error {
	return err.FromError(e).SetCode(configErrCode)
}

// Encoding errors (code 02XXX)
var encodingErrCode = []byte{0x20, 0x00}

// EncodingError are raised when failing to decode/encode a message
func EncodingError(e error) *err.Error {
	return err.FromError(e).SetCode(encodingErrCode)
}

// Data Errors (code 03XXX)
var dataErrCode = []byte{0x30, 0x00}

// DataError is raised when a provided data does not match expected format
func DataError(e error) *err.Error {
	return err.FromError(e).SetCode(dataErrCode)
}

// Solidity Errors (code 031XX)
var solidityErrCode = []byte{0x31, 0x00}

// SolidityError is raised when a data related in transaction crafing is incorrect
func SolidityError(e error) *err.Error {
	return err.FromError(e).SetCode(solidityErrCode)
}

var invalidSigErrCode = []byte{0x31, 0x01}

// InvalidSigError is raised when a solidity method signature is invalid
func InvalidSigError(sig string) *err.Error {
	return err.NewError(fmt.Sprintf("%q is an invalid Solidity method signature (example of valid signature: fctName(address,uint256))", sig)).
		SetCode(invalidSigErrCode)
}
