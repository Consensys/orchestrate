package grpcerror

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Declare custom grpc error code for warning
const Warning codes.Code = 777

// StatusToError convert a status into an internal error
func StatusToError(s *status.Status) *ierror.Error {
	// If an internal error has been detailed then we return it
	for _, detail := range s.Details() {
		if e, ok := detail.(*ierror.Error); ok {
			return e
		}
	}

	// We return an internal error with a error code corresponding to status code
	switch s.Code() {
	case codes.OK:
		return nil
	case Warning:
		return errors.Warningf(s.Message())
	case codes.Canceled:
		return errors.CancelledError(s.Message())
	case codes.Unknown:
		return errors.InternalError(s.Message())
	case codes.InvalidArgument:
		return errors.DataError(s.Message())
	case codes.DeadlineExceeded:
		return errors.DeadlineExceededError(s.Message())
	case codes.NotFound:
		return errors.NotFoundError(s.Message())
	case codes.AlreadyExists:
		return errors.ConstraintViolatedError(s.Message())
	case codes.PermissionDenied:
		return errors.PermissionDeniedError(s.Message())
	case codes.ResourceExhausted:
		return errors.InsufficientResourcesError(s.Message())
	case codes.FailedPrecondition:
		return errors.FailedPreconditionError(s.Message())
	case codes.Aborted:
		return errors.ConflictedError(s.Message())
	case codes.OutOfRange:
		return errors.OutOfRangeError(s.Message())
	case codes.Unimplemented:
		return errors.FeatureNotSupportedError(s.Message())
	case codes.Internal:
		return errors.InternalError(s.Message())
	case codes.Unavailable:
		return errors.GRPCConnectionError(s.Message())
	case codes.DataLoss:
		return errors.DataCorruptedError(s.Message())
	case codes.Unauthenticated:
		return errors.UnauthorizedError(s.Message())
	default:
		return errors.InternalError(s.Message())
	}
}

// ErrorToStatus convert an internal error to a gRPC status
func ErrorToStatus(err error) (s *status.Status) {
	if err == nil {
		return nil
	}

	// If err is built on gRPC status
	se, ok := status.FromError(err)
	if ok {
		// Detail status with internal error generated from se
		s, _ = se.WithDetails(StatusToError(se))
		return
	}

	e := errors.FromError(err)
	switch {
	// Warning
	case errors.IsWarning(e):
		s, _ = status.New(Warning, e.GetMessage()).WithDetails(e)

	// Connection error
	case errors.IsConnectionError(e):
		s, _ = status.New(codes.Unavailable, e.GetMessage()).WithDetails(e)

	// Invalid authentication
	case errors.IsInvalidAuthenticationError(e):
		switch e.GetCode() {
		case errors.Unauthorized:
			s, _ = status.New(codes.Unauthenticated, e.GetMessage()).WithDetails(e)
		case errors.PermissionDenied:
			s, _ = status.New(codes.PermissionDenied, e.GetMessage()).WithDetails(e)
		}

	// Feature not supported
	case errors.IsFeatureNotSupportedError(e):
		s, _ = status.New(codes.Unimplemented, e.GetMessage()).WithDetails(e)

	// Invalid state
	case errors.IsInvalidStateError(e):
		switch {
		case errors.IsFailedPreconditionError(e):
			s, _ = status.New(codes.FailedPrecondition, e.GetMessage()).WithDetails(e)
		case errors.IsConflictedError(e):
			s, _ = status.New(codes.Aborted, e.GetMessage()).WithDetails(e)
		}

	// Data error
	case errors.IsDataError(e):
		switch e.GetCode() {
		case errors.OutOfRange:
			s, _ = status.New(codes.OutOfRange, e.GetMessage()).WithDetails(e)
		default:
			s, _ = status.New(codes.InvalidArgument, e.GetMessage()).WithDetails(e)
		}

	// Insufficient resources
	case errors.IsInsufficientResourcesError(e):
		s, _ = status.New(codes.ResourceExhausted, e.GetMessage()).WithDetails(e)

	// Operator intervention error
	case errors.IsOperatorInterventionError(e):
		switch e.GetCode() {
		case errors.Canceled:
			s, _ = status.New(codes.Canceled, e.GetMessage()).WithDetails(e)
		case errors.DeadlineExceeded:
			s, _ = status.New(codes.DeadlineExceeded, e.GetMessage()).WithDetails(e)
		}

	// Storage error
	case errors.IsStorageError(e):
		switch {
		case errors.IsConstraintViolatedError(e):
			s, _ = status.New(codes.AlreadyExists, e.GetMessage()).WithDetails(e)
		case errors.IsNotFoundError(e):
			s, _ = status.New(codes.NotFound, e.GetMessage()).WithDetails(e)
		}

	// Data corrupted error
	case errors.IsDataCorruptedError(e):
		s, _ = status.New(codes.DataLoss, e.GetMessage()).WithDetails(e)

	// Internal error
	case errors.IsInternalError(e):
		s, _ = status.New(codes.Internal, e.GetMessage()).WithDetails(e)
	}

	if s == nil {
		// Error is unknown
		s, _ = status.New(codes.Unknown, e.GetMessage()).WithDetails(e)
	}

	return
}
