package grpcerror

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor returns a grpc unary interceptor (middleware) that allows
// to intercept translate grpc errors into internal errors
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Invoke
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// Convert error into grpc status and extract internal error
			return StatusToError(status.Convert(err))
		}
		return nil
	}
}

// UnaryClientInterceptor returns a grpc stream interceptor (middleware) that allows
// to intercept translate grpc errors into internal errors
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		stream, err := streamer(ctx, desc, cc, method, opts...)
		estream := &errorClientStream{stream}
		if err != nil {
			return estream, StatusToError(status.Convert(err))
		}
		return estream, nil
	}
}

type errorClientStream struct {
	grpc.ClientStream
}

func (s *errorClientStream) SendMsg(m interface{}) error {
	if err := s.ClientStream.SendMsg(m); err != nil {
		if err == io.EOF {
			return err
		}
		return StatusToError(status.Convert(err))
	}
	return nil
}

func (s *errorClientStream) RecvMsg(m interface{}) error {
	if err := s.ClientStream.RecvMsg(m); err != nil {
		if err == io.EOF {
			return err
		}
		return StatusToError(status.Convert(err))
	}
	return nil
}

func (s *errorClientStream) CloseSend() error {
	if err := s.ClientStream.CloseSend(); err != nil {
		return StatusToError(status.Convert(err))
	}
	return nil
}
