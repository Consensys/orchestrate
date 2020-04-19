package grpcutil

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	Unary        = "unary"
	ClientStream = "client_stream"
	ServerStream = "server_stream"
	BidiStream   = "bidi_stream"
)

var (
	Codes = []codes.Code{
		codes.OK,
		codes.Canceled,
		codes.Unknown,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss,
	}
)

func SplitMethodName(fullMethod string) (service, method string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}
	return "unknown_service", "unknown_method"
}

func TypeFromMethodInfo(mInfo *grpc.MethodInfo) string {
	switch {
	case !mInfo.IsClientStream && !mInfo.IsServerStream:
		return Unary
	case mInfo.IsClientStream && !mInfo.IsServerStream:
		return ClientStream
	case !mInfo.IsClientStream && mInfo.IsServerStream:
		return ServerStream
	default:
		return BidiStream
	}
}

func TypeFromStreamServerInfo(info *grpc.StreamServerInfo) string {
	switch {
	case info.IsClientStream && !info.IsServerStream:
		return ClientStream
	case !info.IsClientStream && info.IsServerStream:
		return ServerStream
	default:
		return BidiStream
	}
}
