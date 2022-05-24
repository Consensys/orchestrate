package errors

import (
	ierror "github.com/consensys/orchestrate/pkg/types/error"
)

// Error codes are uint64 for performances purposes but should be seen as 5 nibbles hex codes
const (
	// Warnings (class 01XXX)
	Warning uint64 = 1 << 12

	// Faucet (class 012XX)
	Faucet = Warning + 2<<8 // Faucet credit denied

	// Invalid Nonce warnings(class 013xx)
	InvalidNonce        = Warning + 3<<8
	InvalidNonceTooHigh = InvalidNonce + 1
	InvalidNonceTooLow  = InvalidNonce + 2

	// Connection Errors (class 08XXX)
	Connection         uint64 = 8 << 12
	KafkaConnection           = Connection + 1<<8 // Kafka Connection error (subclass 081XX)
	HTTPConnection            = Connection + 2<<8 // HTTP Connection error (subclass 082XX)
	EthConnection             = Connection + 3<<8 // Ethereum Connection error (subclass 083XX)
	RedisConnection           = Connection + 5<<8 // Redis Connection error (subclass 085XX)
	PostgresConnection        = Connection + 6<<8 // Postgres Connection error (subclass 086XX)
	ServiceConnection         = Connection + 7<<8 // Service Connection error (subclass 087XX)

	// Authentication Errors (class 09XXX)
	InvalidAuthentication uint64 = 9 << 12
	Unauthorized                 = InvalidAuthentication + 1 // Invalid request credentials (code 09001)
	PermissionDenied             = InvalidAuthentication + 2 // no permission to execute operation (code 09002)

	// Feature Not Supported Errors (class 0AXXX)
	FeatureNotSupported uint64 = 10 << 12

	// Invalid State (class 24XXX)
	InvalidState uint64 = 2<<16 + 4<<12
	Conflicted          = InvalidState + 2<<8 // Conflict with current system state (subclass 242XX)

	// Message Errors (class 42XXX)
	Data               uint64 = 4<<16 + 2<<12
	Encoding                  = Data + 1<<8  // Invalid Encoding (subclass 421XX)
	Solidity                  = Data + 2<<8  // Solidity Errors (subclass 422XX)
	InvalidSignature          = Solidity + 1 // Invalid method/event signature (code 42201)
	InvalidArgsCount          = Solidity + 2 // Invalid arguments count (code 42202)
	InvalidArg                = Solidity + 3 // Invalid argument format (code 42203)
	InvalidTopicsCount        = Solidity + 4 // Invalid count of topics in receipt (code 42204)
	InvalidLog                = Solidity + 5 // Invalid event log (code 42205)
	InvalidFormat             = Data + 3<<8  // Invalid format (subclass 423XX)
	InvalidParameter          = Data + 4<<8  // Invalid parameter provided (subclass 424XX)

	// Ethereum error (class BEXXX)
	Ethereum        uint64 = 11<<16 + 14<<12
	NonceTooLow            = Ethereum + 1
	InvalidNonceErr        = Ethereum + 2

	// Cryptographic operation error (class C0XXX)
	CryptoOperation               uint64 = 12 << 16
	InvalidCryptographicSignature        = CryptoOperation + 1 // Invalid signature during cryptographic verification (subclass C0001)

	// Storage Error (class DBXXX)
	Storage            uint64 = 13<<16 + 11<<12
	ConstraintViolated        = Storage + 1<<8         // Storage constraint violated (subclass DB1XX)
	AlreadyExists             = ConstraintViolated + 1 // A resource with same index already exists (code DB101)
	NotFound                  = Storage + 2<<8         // Not found (subclass DB2XX)

	// Configuration errors (class F0XXX)
	Config uint64 = 15 << 16

	// Internal errors (class FFXXX)
	Internal          uint64 = 15<<16 + 15<<12
	DataCorrupted            = Internal + 1<<8 // Message corrupted (subclass FF1XX)
	DependencyFailure        = Internal + 2<<8 // Message corrupted (subclass FF2XX)
)

// Warningf are raised to indicate a warning
func Warningf(format string, a ...interface{}) *ierror.Error {
	return Errorf(Warning, format, a...)
}

// IsWarning indicate whether an error is a Warning
func IsWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Warning)
}

// FaucetWarning are raised when a faucet credit has been denied
func FaucetWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(Faucet, format, a...)
}

// IsFaucetWarning indicate whether an error is a faucet Warning
func IsFaucetWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Faucet)
}

// InvalidNonceWarning are raised when an invalid nonce is detected
func InvalidNonceWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidNonce, format, a...)
}

// IsInvalidNonceWarning indicate whether an error is an invalid nonce Warning
func IsInvalidNonceWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidNonce)
}

// NonceTooHighWarning are raised when about to send a transaction with nonce too high
func NonceTooHighWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidNonceTooHigh, format, a...)
}

// NonceTooLowWarning are raised when about to send a transaction with nonce too low
func NonceTooLowWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidNonceTooLow, format, a...)
}

// NonceTooLowWarning are raised when about to send a transaction with nonce too low
func InvalidNonceError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidNonceErr, format, a...)
}

// ConnectionError is raised when failing to connect to an external service
func ConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Connection, format, a...)
}

// IsConnectionError indicate whether an error is a Connection error
func IsConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Connection)
}

// KafkaConnectionError is raised when failing to connect to Kafka
func KafkaConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(KafkaConnection, format, a...)
}

// IsKafkaConnectionError indicate whether an error is a KafkaConnection error
func IsKafkaConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), KafkaConnection)
}

// HTTPConnectionError is raised when failing to connect over http
func HTTPConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(HTTPConnection, format, a...)
}

// EthConnectionError is raised when failing to connect to Ethereum client jsonRPC API
func EthConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(EthConnection, format, a...)
}

// IsEthConnectionError indicate whether an error is a EthConnection error
func IsEthConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), EthConnection)
}

// RedisConnectionError is raised when failing to connect to Redis
func RedisConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(RedisConnection, format, a...)
}

// PostgresConnectionError is raised when failing to connect to Postgres
func PostgresConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(PostgresConnection, format, a...)
}

// IsPostgresConnectionError indicate whether an error is a Postgres connection error
func IsPostgresConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), PostgresConnection)
}

// ServiceConnectionError is raised when failing to connect to another service
func ServiceConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(ServiceConnection, format, a...)
}

// IsServiceConnectionError indicate whether an error is a Service connection error
func IsServiceConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), ServiceConnection)
}

// InvalidAuthenticationError is raised when access to an operation has been denied
func InvalidAuthenticationError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidAuthentication, format, a...)
}

// AuthenticationError indicate whether an error is an authentication error
func IsInvalidAuthenticationError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidAuthentication)
}

// UnauthorizedError is raised when authentication credentials are invalid
func UnauthorizedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Unauthorized, format, a...)
}

func IsUnauthorizedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Unauthorized)
}

// PermissionDeniedError is raised when authentication credentials are invalid
func PermissionDeniedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(PermissionDenied, format, a...)
}

// FeatureNotSupportedError is raised when using a feature which is not implemented
func FeatureNotSupportedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(FeatureNotSupported, format, a...)
}

// IsFeatureNotSupportedError indicate whether an error is a feature not supported error
func IsFeatureNotSupportedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), FeatureNotSupported)
}

// InvalidStateError is raised when system state blocks operation execution
func InvalidStateError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidState, format, a...)
}

// AuthenticationError indicate whether an error is an invalid state error
func IsInvalidStateError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidState)
}

// ConflictedError is raised when operation could not be completed due to a
// conflict with the current state of the target resource
//
// User might be able to resolve the conflict and resubmit operation
func ConflictedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Conflicted, format, a...)
}

// IsConflictedError indicate whether an error is an conflicted error
func IsConflictedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Conflicted)
}

// DataError is raised when a provided Message does not match expected format
func DataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Data, format, a...)
}

// IsDataError indicate whether an error is a Message error
func IsDataError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Data)
}

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Encoding, format, a...)
}

// IsEncodingError indicate whether an error is a EncodingError error
func IsEncodingError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Encoding)
}

// SolidityError is raised when a Message related in transaction crafting is incorrect
func SolidityError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Solidity, format, a...)
}

// IsSolidityError indicate whether an error is a Solidity error
func IsSolidityError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Solidity)
}

// InvalidSignatureError is raised when a Solidity method signature is invalid
func InvalidSignatureError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidSignature, format, a...)
}

// InvalidArgsCountError is raised when invalid arguments count is provided to craft a transaction
func InvalidArgsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidArgsCount, format, a...)
}

// InvalidArgError is raised when invalid argument is provided to craft a transaction
func InvalidArgError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidArg, format, a...)
}

func IsInvalidArgError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidArg)
}

// InvalidTopicsCountError is raised when topics count is in receipt
func InvalidTopicsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidTopicsCount, format, a...)
}

// InvalidEventDataError is raised when event Message is invalid
func InvalidEventDataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidLog, format, a...)
}

// InvalidFormatError is raised when a Message does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidFormat, format, a...)
}

// IsInvalidFormatError indicate whether an error is an invalid format error
func IsInvalidFormatError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidFormat)
}

// InvalidParameterError is raised when a provided parameter invalid
func InvalidParameterError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidParameter, format, a...)
}

// IsInvalidParameterError indicate whether an error is an invalid parameter error
func IsInvalidParameterError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidParameter)
}

// EthereumError is raised when JSON-RPC call returns an error (such as Nonce too Low)
func EthereumError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Ethereum, format, a...)
}

// IsEthereumError indicate whether an error is an Ethereum error
func IsEthereumError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Ethereum)
}

// NonceTooLowError is raised when JSON-RPC returns a "Nonce too low" when sending a transaction
func NonceTooLowError(format string, a ...interface{}) *ierror.Error {
	return Errorf(NonceTooLow, format, a...)
}

// CryptoOperationError is raised when failing a cryptographic operation
func CryptoOperationError(format string, a ...interface{}) *ierror.Error {
	return Errorf(CryptoOperation, format, a...)
}

// IsCryptoOperationError indicate whether an error is a cryptographic operation error
func IsCryptoOperationError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), CryptoOperation)
}

// InvalidCryptographicSignature is raised when failing a signature cryptographic verification
func InvalidCryptographicSignatureError(format string, a ...interface{}) *ierror.Error {
	return Errorf(InvalidCryptographicSignature, format, a...)
}

// IsInvalidCryptographicSignatureError indicate whether an error is a signature cryptographic verification error
func IsInvalidCryptographicSignatureError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidCryptographicSignature)
}

// InvalidCryptographicSignature

// StorageError is raised when an error is encountered while accessing stored Message
func StorageError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Storage, format, a...)
}

// IsStorageError indicate whether an error is a Storage error
func IsStorageError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Storage)
}

// ConstraintViolatedError is raised when a Message constraint has been violated
func ConstraintViolatedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(ConstraintViolated, format, a...)
}

// IsConstraintViolatedError indicate whether an error is a constraint violated error
func IsConstraintViolatedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), ConstraintViolated)
}

// AlreadyExistsError is raised when a Message constraint has been violated
func AlreadyExistsError(format string, a ...interface{}) *ierror.Error {
	return Errorf(AlreadyExists, format, a...)
}

// IsAlreadyExistsError indicate whether an error is an already exists error
func IsAlreadyExistsError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AlreadyExists)
}

// NoDataFoundError is raised when accessing a missing Message
func NotFoundError(format string, a ...interface{}) *ierror.Error {
	return Errorf(NotFound, format, a...)
}

// IsNotFoundError indicate whether an error is a no Message found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotFound)
}

// ConfigError is raised when an error is encountered while loading configuration
func ConfigError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Config, format, a...)
}

// InternalError is raised when an unknown exception is met
func InternalError(format string, a ...interface{}) *ierror.Error {
	return Errorf(Internal, format, a...)
}

// IsInternalError indicate whether an error is an Internal error
func IsInternalError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Internal)
}

func DependencyFailureError(format string, a ...interface{}) *ierror.Error {
	return Errorf(DependencyFailure, format, a...)
}

func IsDependencyFailureError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), DependencyFailure)
}

// DataCorruptedError is raised loading a corrupted Message
func DataCorruptedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(DataCorrupted, format, a...)
}

// IsDataCorruptedError indicate whether an error is a data corrupted error
func IsDataCorruptedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), DataCorrupted)
}
