package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarningf(t *testing.T) {
	e := Warningf("test")
	assert.Equal(t, uint64(4096), e.GetCode(), "Warning code should be correct")
	assert.Equal(t, "01000", e.Hex(), "Hex representation should be correct")
}

func TestRetryWarning(t *testing.T) {
	e := RetryWarning("test")
	assert.Equal(t, uint64(4352), e.GetCode(), "RetryWarning code should be correct")
	assert.True(t, IsWarning(e), "RetryWarning should be a connection error")
	assert.Equal(t, "01100", e.Hex(), "RetryWarning Hex representation should be correct")
}

func TestFaucetWarning(t *testing.T) {
	e := FaucetWarning("test")
	assert.Equal(t, uint64(4608), e.GetCode(), "FaucetWarning code should be correct")
	assert.True(t, IsWarning(e), "FaucetWarning should be a connection error")
	assert.True(t, IsFaucetWarning(e), "FaucetWarning should be a connection error")
	assert.Equal(t, "01200", e.Hex(), "FaucetWarning Hex representation should be correct")
}

func TestInvalidNonceWarning(t *testing.T) {
	e := InvalidNonceWarning("test")
	assert.Equal(t, uint64(4864), e.GetCode(), "InvalidNonceWarning code should be correct")
	assert.True(t, IsWarning(e), "InvalidNonceWarning should be a connection error")
	assert.True(t, IsInvalidNonceWarning(e), "InvalidNonceWarning should be a connection error")
	assert.Equal(t, "01300", e.Hex(), "InvalidNonceWarning Hex representation should be correct")
}

func TestNonceTooHighWarning(t *testing.T) {
	e := NonceTooHighWarning("test")
	assert.Equal(t, uint64(4865), e.GetCode(), "NonceTooHighWarning code should be correct")
	assert.True(t, IsWarning(e), "NonceTooHighWarning should be a connection error")
	assert.True(t, IsInvalidNonceWarning(e), "NonceTooHighWarning should be a connection error")
	assert.Equal(t, "01301", e.Hex(), "NonceTooHighWarning Hex representation should be correct")
}

func TestNonceTooLowWarning(t *testing.T) {
	e := NonceTooLowWarning("test")
	assert.Equal(t, uint64(4866), e.GetCode(), "NonceTooLowWarning code should be correct")
	assert.True(t, IsWarning(e), "NonceTooLowWarning should be a connection error")
	assert.True(t, IsInvalidNonceWarning(e), "NonceTooLowWarning should be a connection error")
	assert.Equal(t, "01302", e.Hex(), "NonceTooLowWarning Hex representation should be correct")
}

func TestConnectionError(t *testing.T) {
	e := ConnectionError("test")
	assert.Equal(t, uint64(32768), e.GetCode(), "ConnectionError code should be correct")
	assert.Equal(t, "08000", e.Hex(), "Hex representation should be correct")
}

func TestKafkaConnectionError(t *testing.T) {
	e := KafkaConnectionError("test")
	assert.Equal(t, uint64(33024), e.GetCode(), "KafkaConnectionError code should be correct")
	assert.True(t, IsConnectionError(e), "KafkaConnectionError should be a connection error")
	assert.Equal(t, "08100", e.Hex(), "KafkaConnectionError Hex representation should be correct")
}

func TestHTTPConnectionError(t *testing.T) {
	e := HTTPConnectionError("test")
	assert.Equal(t, uint64(33280), e.GetCode(), "HTTPConnectionError code should be correct")
	assert.True(t, IsConnectionError(e), "HTTPConnectionError should be a connection error")
	assert.Equal(t, "08200", e.Hex(), "HTTPConnectionError Hex representation should be correct")
}

func TestEthConnectionError(t *testing.T) {
	e := EthConnectionError("test")
	assert.Equal(t, uint64(33536), e.GetCode(), "EthConnectionError code should be correct")
	assert.True(t, IsConnectionError(e), "EthConnectionError should be a connection error")
	assert.Equal(t, "08300", e.Hex(), "EthConnectionError Hex representation should be correct")
}

func TestGRPCConnectionError(t *testing.T) {
	e := GRPCConnectionError("test")
	assert.Equal(t, uint64(33792), e.GetCode(), "GRPCConnectionError code should be correct")
	assert.True(t, IsConnectionError(e), "GRPCConnectionError should be a connection error")
	assert.Equal(t, "08400", e.Hex(), "GRPCConnectionError Hex representation should be correct")
}

func TestRedisConnectionError(t *testing.T) {
	e := RedisConnectionError("test")
	assert.Equal(t, uint64(34048), e.GetCode(), "RedisConnectionError code should be correct")
	assert.True(t, IsConnectionError(e), "RedisConnectionError should be a connection error")
	assert.Equal(t, "08500", e.Hex(), "RedisConnectionError Hex representation should be correct")
}

func TestInvalidAuthenticationError(t *testing.T) {
	e := InvalidAuthenticationError("test")
	assert.Equal(t, uint64(36864), e.GetCode(), "InvalidAuthenticationError code should be correct")
	assert.Equal(t, "09000", e.Hex(), "Hex representation should be correct")
}

func TestUnauthenticatedError(t *testing.T) {
	e := UnauthenticatedError("test")
	assert.Equal(t, uint64(36865), e.GetCode(), "UnauthenticatedError code should be correct")
	assert.True(t, IsInvalidAuthenticationError(e), "UnauthenticatedError should be a connection error")
	assert.Equal(t, "09001", e.Hex(), "UnauthenticatedError Hex representation should be correct")
}

func TestPermissionDeniedError(t *testing.T) {
	e := PermissionDeniedError("test")
	assert.Equal(t, uint64(36866), e.GetCode(), "PermissionDeniedError code should be correct")
	assert.True(t, IsInvalidAuthenticationError(e), "PermissionDeniedError should be a connection error")
	assert.Equal(t, "09002", e.Hex(), "PermissionDeniedError Hex representation should be correct")
}

func TestConfigError(t *testing.T) {
	e := ConfigError("test")
	assert.Equal(t, uint64(983040), e.GetCode(), "ConfigError code should be correct")
	assert.Equal(t, "F0000", e.Hex(), "Hex representation should be correct")
}

func TestFeatureNotSupportedError(t *testing.T) {
	e := FeatureNotSupportedError("test")
	assert.Equal(t, uint64(40960), e.GetCode(), "FeatureNotSupportedError code should be correct")
	assert.Equal(t, "0A000", e.Hex(), "Hex representation should be correct")
}

func TestInvalidStateError(t *testing.T) {
	e := InvalidStateError("test")
	assert.Equal(t, uint64(147456), e.GetCode(), "InvalidStateError code should be correct")
	assert.Equal(t, "24000", e.Hex(), "Hex representation should be correct")
}

func TestFailedPreconditionError(t *testing.T) {
	e := FailedPreconditionError("test")
	assert.Equal(t, uint64(147712), e.GetCode(), "FailedPreconditionError code should be correct")
	assert.True(t, IsInvalidStateError(e), "FailedPreconditionError should be InvalidStateError")
	assert.Equal(t, "24100", e.Hex(), "FailedPreconditionError Hex representation should be correct")
}

func TestConflictedError(t *testing.T) {
	e := ConflictedError("test")
	assert.Equal(t, uint64(147968), e.GetCode(), "ConflictedError code should be correct")
	assert.True(t, IsInvalidStateError(e), "ConflictedError should be InvalidStateError")
	assert.Equal(t, "24200", e.Hex(), "ConflictedError Hex representation should be correct")
}

func TestDataError(t *testing.T) {
	e := DataError("test")
	assert.Equal(t, uint64(270336), e.GetCode(), "DataError code should be correct")
	assert.Equal(t, "42000", e.Hex(), "Hex representation should be correct")
}

func TestOutOfRange(t *testing.T) {
	e := OutOfRangeError("test")
	assert.Equal(t, uint64(270337), e.GetCode(), "OutOfRangeError code should be correct")
	assert.Equal(t, "42001", e.Hex(), "Hex representation should be correct")
}

func TestEncodingError(t *testing.T) {
	e := EncodingError("test")
	assert.Equal(t, uint64(270592), e.GetCode(), "EncodingError code should be correct")
	assert.True(t, IsDataError(e), "EncodingError should be a data error")
	assert.Equal(t, "42100", e.Hex(), "Hex representation should be correct")
}

func TestSolidityError(t *testing.T) {
	e := SolidityError("test")
	assert.Equal(t, uint64(270848), e.GetCode(), "SolidityError code should be correct")
	assert.True(t, IsDataError(e), "SolidityError should be a data error")
	assert.Equal(t, "42200", e.Hex(), "Hex representation should be correct")
}

func TestInvalidSignatureError(t *testing.T) {
	e := InvalidSignatureError("test")
	assert.Equal(t, uint64(270849), e.GetCode(), "InvalidSignatureError code should be correct")
	assert.True(t, IsDataError(e), "InvalidSignatureError should be a data error")
	assert.True(t, IsSolidityError(e), "IsSolidityError should be a data error")
	assert.Equal(t, "42201", e.Hex(), "Hex representation should be correct")
}

func TestInvalidInvalidArgsCountError(t *testing.T) {
	e := InvalidArgsCountError("test")
	assert.Equal(t, uint64(270850), e.GetCode(), "InvalidArgsCountError code should be correct")
	assert.True(t, IsDataError(e), "InvalidArgsCountError should be a data error")
	assert.True(t, IsSolidityError(e), "InvalidArgsCountError should be a data error")
	assert.Equal(t, "42202", e.Hex(), "Hex representation should be correct")
}

func TestInvalidArgError(t *testing.T) {
	e := InvalidArgError("test")
	assert.Equal(t, uint64(270851), e.GetCode(), "InvalidArgCode code should be correct")
	assert.True(t, IsDataError(e), "InvalidArgCode should be a data error")
	assert.True(t, IsSolidityError(e), "InvalidArgCode should be a data error")
	assert.Equal(t, "42203", e.Hex(), "Hex representation should be correct")
}

func TestInvalidTopicsCountError(t *testing.T) {
	e := InvalidTopicsCountError("test")
	assert.Equal(t, uint64(270852), e.GetCode(), "InvalidTopicsCountError code should be correct")
	assert.True(t, IsDataError(e), "InvalidTopicsCountError should be a data error")
	assert.True(t, IsSolidityError(e), "InvalidTopicsCountError should be a data error")
	assert.Equal(t, "42204", e.Hex(), "Hex representation should be correct")
}

func TestInvalidEventDataError(t *testing.T) {
	e := InvalidEventDataError("test")
	assert.Equal(t, uint64(270853), e.GetCode(), "InvalidEventDataError code should be correct")
	assert.True(t, IsDataError(e), "InvalidEventDataError should be a data error")
	assert.True(t, IsSolidityError(e), "InvalidEventDataError should be a data error")
	assert.Equal(t, "42205", e.Hex(), "Hex representation should be correct")
}

func TestInvalidFormatError(t *testing.T) {
	e := InvalidFormatError("test")
	assert.Equal(t, uint64(271104), e.GetCode(), "InvalidFormatError code should be correct")
	assert.True(t, IsDataError(e), "InvalidFormatError should be a data error")
	assert.Equal(t, "42300", e.Hex(), "Hex representation should be correct")
}

func TestInvalidParameterError(t *testing.T) {
	e := InvalidParameterError("test")
	assert.Equal(t, uint64(271360), e.GetCode(), "InvalidParameterError code should be correct")
	assert.True(t, IsDataError(e), "InvalidParameterError should be a data error")
	assert.True(t, IsInvalidParameterError(e), "InvalidParameterError should be a InvalidParameterError")
	assert.Equal(t, "42400", e.Hex(), "Hex representation should be correct")
}

func TestInsufficientResourcesError(t *testing.T) {
	e := InsufficientResourcesError("test")
	assert.Equal(t, uint64(339968), e.GetCode(), "InsufficientResourcesError code should be correct")
	assert.Equal(t, "53000", e.Hex(), "Hex representation should be correct")
}

func TestOperatorInterventionError(t *testing.T) {
	e := OperatorInterventionError("test")
	assert.Equal(t, uint64(356352), e.GetCode(), "OperatorInterventionError code should be correct")
	assert.Equal(t, "57000", e.Hex(), "Hex representation should be correct")
}

func TestCancelledError(t *testing.T) {
	e := CancelledError("test")
	assert.Equal(t, uint64(356353), e.GetCode(), "CancelledError code should be correct")
	assert.True(t, IsOperatorInterventionError(e), "CancelledError should be a OperatorInterventionError")
	assert.Equal(t, "57001", e.Hex(), "Hex representation should be correct")
}

func TestDeadlineExceededError(t *testing.T) {
	e := DeadlineExceededError("test")
	assert.Equal(t, uint64(356354), e.GetCode(), "DeadlineExceededError code should be correct")
	assert.True(t, IsOperatorInterventionError(e), "DeadlineExceededError should be a OperatorInterventionError")
	assert.Equal(t, "57002", e.Hex(), "Hex representation should be correct")
}

func TestCryptoOperationError(t *testing.T) {
	e := CryptoOperationError("test")
	assert.Equal(t, uint64(786432), e.GetCode(), "CryptoOperationError code should be correct")
	assert.True(t, IsCryptoOperationError(e), "CryptoOperationError should be a CryptoOperationError")
	assert.Equal(t, "C0000", e.Hex(), "Hex representation should be correct")
}

func TestStorageError(t *testing.T) {
	e := StorageError("test")
	assert.Equal(t, uint64(897024), e.GetCode(), "StorageError code should be correct")
	assert.Equal(t, "DB000", e.Hex(), "Hex representation should be correct")
}

func TestConstraintViolatedError(t *testing.T) {
	e := ConstraintViolatedError("test")
	assert.Equal(t, uint64(897280), e.GetCode(), "ConstraintViolatedError code should be correct")
	assert.True(t, IsStorageError(e), "ConstraintViolatedError should be a StorageError")
	assert.Equal(t, "DB100", e.Hex(), "Hex representation should be correct")
}

func TestAlreadyExistsError(t *testing.T) {
	e := AlreadyExistsError("test")
	assert.Equal(t, uint64(897281), e.GetCode(), "AlreadyExistsError code should be correct")
	assert.True(t, IsStorageError(e), "AlreadyExistsError should be StorageError")
	assert.True(t, IsConstraintViolatedError(e), "AlreadyExistsError should be a ConstraintViolatedError")
	assert.Equal(t, "DB101", e.Hex(), "Hex representation should be correct")
}

func TestNotFoundError(t *testing.T) {
	e := NotFoundError("test")
	assert.Equal(t, uint64(897536), e.GetCode(), "NotFoundError code should be correct")
	assert.True(t, IsStorageError(e), "NotFoundError should be a data error")
	assert.True(t, IsNotFoundError(e), "DataCorruptedError should be a data error")
	assert.Equal(t, "DB200", e.Hex(), "Hex representation should be correct")
}

func TestInternalError(t *testing.T) {
	e := InternalError("test")
	assert.Equal(t, uint64(1044480), e.GetCode(), "InternalError code should be correct")
	assert.Equal(t, "FF000", e.Hex(), "Hex representation should be correct")
}

func TestDataCorruptedError(t *testing.T) {
	e := DataCorruptedError("test")
	assert.Equal(t, uint64(1044736), e.GetCode(), "DataCorruptedError code should be correct")
	assert.True(t, IsInternalError(e), "DataCorruptedError should be an InternalError")
	assert.Equal(t, "FF100", e.Hex(), "Hex representation should be correct")
}
