package grpcerror

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a grpc unary server interceptor (middleware) that allows
// to intercept internal errors
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, ErrorToStatus(err).Err()
		}
		return resp, nil
	}
}

// StreamServerInterceptor returns a grpc streaming server interceptor for panic recovery.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		estream := &errorServerStream{stream}
		err := handler(srv, estream)
		if err != nil {
			return ErrorToStatus(err).Err()
		}
		return nil
	}
}

type errorServerStream struct {
	grpc.ServerStream
}

func (s *errorServerStream) SetHeader(md metadata.MD) error {
	if err := s.ServerStream.SetHeader(md); err != nil {
		return StatusToError(status.Convert(err))
	}
	return nil
}

func (s *errorServerStream) SendHeader(md metadata.MD) error {
	if err := s.ServerStream.SendHeader(md); err != nil {
		return StatusToError(status.Convert(err))
	}
	return nil
}

func (s *errorServerStream) SendMsg(m interface{}) error {
	if err := s.ServerStream.SendMsg(m); err != nil {
		return StatusToError(status.Convert(err))
	}
	return nil
}

func (s *errorServerStream) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		if err == io.EOF {
			return err
		}
		return StatusToError(status.Convert(err))
	}
	return nil
}
