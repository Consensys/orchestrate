package errors

import (
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Error codes are uint64 for perfomances purposes but should be seen as 5 nibbles hex codes
const (
	// Warnings (class 01XXX)
	warning       = 1 << 12
	retryWarning  = warning + 1<<8 // Retries (subclass 011XX)
	faucetWarning = warning + 2<<8 // Faucet credit denied (subclass 012XX)

	// Connnection Errors (class 08XXX)
	connection      = 8 << 12
	kafkaConnection = connection + 1<<8 // Kafka connection error (subclass 081XX)
	httpConnection  = connection + 2<<8 // HTTP connection error (subclass 082XX)
	ethConnection   = connection + 3<<8 // Ethereum connection error (subclass 083XX)
	grpcConnection  = connection + 4<<8 // GRPC connection error (subclass 084XX)

	// Authentication Errors (class 09XXX)
	invalidAuthentication = 9 << 12
	unauthenticated       = invalidAuthentication + 1 // Invalid request credentials (code 09001)
	permissionDenied      = invalidAuthentication + 2 // no permission to execute operation (code 09002)

	// Feature Not Supported Errors (class 0AXXX)
	featureNotSupported = 10 << 12

	// Data Errors (class 42XXX)
	data               = 4<<16 + 2<<12
	outOfRange         = data + 1     // Out of range (code 42001)
	encoding           = data + 1<<8  // Invalid encoding (subclass 421XX)
	solidity           = data + 2<<8  // Solidity Errors (subclass 422XX)
	invalidSignature   = solidity + 1 // Invalid method/event signature (code 42201)
	invalidArgsCount   = solidity + 2 // Invalid arguments count (code 42202)
	invalidArg         = solidity + 3 // Invalid argument format (code 42203)
	invalidTopicsCount = solidity + 4 // Invalid count of topics in receipt (code 42204)
	invalidLog         = solidity + 5 // Invalid event log (code 42205)
	invalidFormat      = data + 3<<8  // Invalid format (subclass 423XX)

	// Insuficient resources (class 53XXX)
	insufficientResources = 5<<16 + 3<<12

	// Operation intervention error (class 57XXX)
	operatorIntervention = 5<<16 + 7<<12 //
	canceled             = operatorIntervention + 1
	deadlineExceeded     = operatorIntervention + 2

	// Storage Error (class DBXXX)
	storage            = 13<<16 + 11<<12
	constraintViolated = storage + 1<<8 // Storage constraint violated (subclass DB1XX)
	notFound           = storage + 2<<8 // Not found (subclass DB2XX)
	dataCorrupted      = storage + 3<<8 // Data corrupted (subclass DB3XX)

	// Configuration errors (class F0XXX)
	config = 15 << 16

	// Internal errors (class FFXXX)
	internal = 15<<16 + 15<<12
)

// Warning are raised to indicate
func Warning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(warning)
}

// IsWarning indicate whether an error is a warning
func IsWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), warning)
}

// RetryWarning are raised when failing to connect to a service and retrying
func RetryWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(retryWarning)
}

// FaucetWarning are raised when a faucet credit has been denied
func FaucetWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(faucetWarning)
}

// IsFaucetWarning indicate whether an error is a faucet warning
func IsFaucetWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), faucetWarning)
}

// ConnectionError is raised when failing to connect to an external service
func ConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(connection)
}

// IsConnectionError indicate whether an error is a connection error
func IsConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), connection)
}

// KafkaConnectionError is raised when failing to connect to Kafka
func KafkaConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(kafkaConnection)
}

// HTTPConnectionError is raised when failing to connect over http
func HTTPConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(httpConnection)
}

// EthConnectionError is raised when failing to connect to Ethereum client jsonRPC API
func EthConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(ethConnection)
}

// GRPCConnectionError is raised when failing to connect to a GRPC server
func GRPCConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(grpcConnection)
}

// InvalidAuthenticationError is raised when access to an operation has been denied
func InvalidAuthenticationError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidAuthentication)
}

// AuthenticationError indicate whether an error is an authentication error
func IsInvalidAuthenticationError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), invalidAuthentication)
}

// UnauthenticatedError is raised when authentication credentials are invalid
func UnauthenticatedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(unauthenticated)
}

// PermissionDeniedError is raised when authentication credentials are invalid
func PermissionDeniedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(permissionDenied)
}

// FeatureNotSupportedError is raised when using a feature which is not implemented
func FeatureNotSupportedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(featureNotSupported)
}

// DataError is raised when a provided data does not match expected format
func DataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(data)
}

// IsDataError indicate whether an error is a data error
func IsDataError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), data)
}

// OutOfRangeError are raised when an operation was attempted past the valid range
func OutOfRangeError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(outOfRange)
}

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(encoding)
}

// SolidityError is raised when a data related in transaction crafing is incorrect
func SolidityError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(solidity)
}

// IsSolidityError indicate whether an error is a Solidity error
func IsSolidityError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), solidity)
}

// InvalidSignatureError is raised when a solidity method signature is invalid
func InvalidSignatureError(sig string) *ierror.Error {
	return Errorf("%q is an invalid Solidity method signature (example of valid signature: transfer(address,uint256))", sig).
		SetCode(invalidSignature)
}

// InvalidArgsCountError is raised when invalid arguments count is provided to craft a transaction
func InvalidArgsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArgsCount)
}

// InvalidArgError is raised when invalid argument is provided to craft a transaction
func InvalidArgError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArg)
}

// InvalidTopicsCountError is raised when topics count is in receipt
func InvalidTopicsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidTopicsCount)
}

// InvalidEventDataError is raised when event data is invalid
func InvalidEventDataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidLog)
}

// InvalidFormatError is raised when a data does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidFormat)
}

// InsuficientResourcesError is raised when a system can not handle more operations
func InsuficientResourcesError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(insufficientResources)
}

// OperatorInterventionError is raised when an error resulted from an operator interfering with the system
func OperatorInterventionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(operatorIntervention)
}

// IsOperatorInterventionError indicate whether an error is a operator intervention error
func IsOperatorInterventionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), operatorIntervention)
}

// CancelledError is raised when canceling an operation
func CancelledError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(canceled)
}

// DeadlineExceededError is raised when deadline expired before operation could complete
func DeadlineExceededError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(deadlineExceeded)
}

// StorageError is raised when an error is encountered while accessing stored data
func StorageError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(storage)
}

// IsStorageError indicate whether an error is a storage error
func IsStorageError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), storage)
}

// ConstraintViolatedError is raised when a data constraint has been violated
func ConstraintViolatedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(constraintViolated)
}

// IsConstraintViolatedError indicate whether an error is a constraint violated error
func IsConstraintViolatedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), constraintViolated)
}

// NoDataFoundError is raised when accessing a missing data
func NotFoundError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(notFound)
}

// IsNotFoundError indicate whether an error is a no data found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), notFound)
}

// DataCorruptedError is raised loading a corrupted data
func DataCorruptedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(dataCorrupted)
}

// ConfigError is raised when an error is encountered while loading configuration
func ConfigError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(config)
}

// InternalError is raised when an unknown exception is met
func InternalError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...)
}
