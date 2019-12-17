package grpcerror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testStatusToError(t *testing.T, code codes.Code, errcode uint64) {
	err := errors.FromError(StatusToError(status.New(code, "test")))
	assert.NotNil(t, err, "Should error")
	assert.Equal(t, "test", err.GetMessage(), "Error message should be correct")
	assert.Equal(t, errcode, err.GetCode(), "Error code should be correct")
}

func TestStatusToError(t *testing.T) {
	// When status code OK should return nil
	s := status.New(codes.OK, "test")
	e := errors.FromError(StatusToError(s))
	assert.Nil(t, e, "Success should not transform to error")

	// When status code is not OK should return internal error with corresponding error code
	testStatusToError(t, Warning, errors.Warning)
	testStatusToError(t, codes.Canceled, errors.Canceled)
	testStatusToError(t, codes.Unknown, errors.Internal)
	testStatusToError(t, codes.InvalidArgument, errors.Data)
	testStatusToError(t, codes.DeadlineExceeded, errors.DeadlineExceeded)
	testStatusToError(t, codes.NotFound, errors.NotFound)
	testStatusToError(t, codes.AlreadyExists, errors.ConstraintViolated)
	testStatusToError(t, codes.PermissionDenied, errors.PermissionDenied)
	testStatusToError(t, codes.ResourceExhausted, errors.InsufficientResources)
	testStatusToError(t, codes.FailedPrecondition, errors.FailedPrecondition)
	testStatusToError(t, codes.Aborted, errors.Conflicted)
	testStatusToError(t, codes.OutOfRange, errors.OutOfRange)
	testStatusToError(t, codes.Unimplemented, errors.FeatureNotSupported)
	testStatusToError(t, codes.Internal, errors.Internal)
	testStatusToError(t, codes.Unavailable, errors.GRPCConnection)
	testStatusToError(t, codes.DataLoss, errors.DataCorrupted)
	testStatusToError(t, codes.Unauthenticated, errors.Unauthorized)
}

func testErrorToStatus(t *testing.T, err error, code codes.Code) {
	s := ErrorToStatus(err)
	assert.NotNil(t, s, "Status should not be nil error")
	assert.Equal(t, code, s.Code(), "Status code should be correct")

	e := errors.FromError(err)
	assert.Equal(t, e.GetMessage(), s.Message(), "Error message should be correct")

	ee := errors.FromError(StatusToError(s))
	assert.Equal(t, e.GetCode(), ee.GetCode(), "StatusToError should inverse ErrorToStatus, Error code")
	assert.Equal(t, e.GetMessage(), ee.GetMessage(), "StatusToError should inverse ErrorToStatus, Message")
	assert.Equal(t, e.GetComponent(), ee.GetComponent(), "StatusToError should inverse ErrorToStatus, Component")
}

func TestErrorToStatus(t *testing.T) {
	// When error is nil should return nil
	assert.Nil(t, ErrorToStatus(nil), "Status of nil error should be nil")

	// When error is not nil should map to status with expected error code
	testErrorToStatus(t, errors.FaucetWarning("test"), Warning)
	testErrorToStatus(t, errors.KafkaConnectionError("test"), codes.Unavailable)
	testErrorToStatus(t, errors.UnauthorizedError("test"), codes.Unauthenticated)
	testErrorToStatus(t, errors.PermissionDeniedError("test"), codes.PermissionDenied)
	testErrorToStatus(t, errors.FeatureNotSupportedError("test"), codes.Unimplemented)
	testErrorToStatus(t, errors.FailedPreconditionError("test"), codes.FailedPrecondition)
	testErrorToStatus(t, errors.ConflictedError("test"), codes.Aborted)
	testErrorToStatus(t, errors.OutOfRangeError("test"), codes.OutOfRange)
	testErrorToStatus(t, errors.InvalidSignatureError("test"), codes.InvalidArgument)
	testErrorToStatus(t, errors.InsufficientResourcesError("test"), codes.ResourceExhausted)
	testErrorToStatus(t, errors.CancelledError("test"), codes.Canceled)
	testErrorToStatus(t, errors.DeadlineExceededError("test"), codes.DeadlineExceeded)
	testErrorToStatus(t, errors.NotFoundError("test"), codes.NotFound)
	testErrorToStatus(t, errors.ConstraintViolatedError("test"), codes.AlreadyExists)
	testErrorToStatus(t, errors.DataCorruptedError("test"), codes.DataLoss)
	testErrorToStatus(t, errors.Errorf(errors.Internal, "test"), codes.Internal)
	testErrorToStatus(t, fmt.Errorf("test"), codes.Internal)
}
