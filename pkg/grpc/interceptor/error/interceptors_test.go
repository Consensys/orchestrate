// +build unit

package grpcerror

import (
	"context"
	"fmt"
	"io"
	"testing"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpctesting "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testInternal = "internal"
	testFmt      = "fmt"
)

func handlePing(ping *testproto.PingRequest) error {
	if ping.GetErrorCodeReturned() == 0 {
		return nil
	}

	if ping.Value == testInternal {
		return errors.Errorf(256, "test-internal").SetComponent("ping")
	}

	if ping.Value == testFmt {
		return fmt.Errorf("test-fmt")
	}

	return status.New(codes.Code(ping.GetErrorCodeReturned()), ping.GetValue()).Err()
}

type errorPingService struct {
	*grpctesting.TestPingService
}

func (s *errorPingService) PingError(ctx context.Context, ping *testproto.PingRequest) (*testproto.Empty, error) {
	return &testproto.Empty{}, handlePing(ping)
}

func (s *errorPingService) PingList(ping *testproto.PingRequest, stream testproto.TestService_PingListServer) error {
	err := handlePing(ping)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if err := stream.Send(&testproto.PingResponse{Value: ping.Value, Counter: int32(i)}); err != nil {
			return err
		}
	}

	return nil
}

func (s *errorPingService) PingStream(stream testproto.TestService_PingStreamServer) error {
	count := int32(0)
	log.Info("PingStream")
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		err = handlePing(ping)
		if err != nil {
			return err
		}

		if err := stream.Send(&testproto.PingResponse{Value: ping.Value, Counter: count}); err != nil {
			return err
		}
		count++
	}
	return nil
}

type errorInterceptorsSuite struct {
	*grpctesting.InterceptorTestSuite
}

func (s *errorInterceptorsSuite) TestPingErrorSuccess() {
	// Success
	_, err := s.Client.PingError(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "success"},
	)
	assert.Nil(s.T(), err, "Error should be nil on success")
}

func (s *errorInterceptorsSuite) TestPingErrorInternal() {
	// Internal error
	_, err := s.Client.PingError(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "internal", ErrorCodeReturned: 1},
	)
	e := errors.FromError(err)
	assert.Equal(s.T(), "ping", e.GetComponent(), "Error component should have been set")
	assert.Equal(s.T(), "00100", e.Hex(), "Error code should be correct")
	assert.Equal(s.T(), "test-internal", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingErrorGRPC() {
	// gRPC error
	_, err := s.Client.PingError(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "test-error", ErrorCodeReturned: uint32(codes.OutOfRange)},
	)

	e := errors.FromError(err)
	assert.Equal(s.T(), errors.OutOfRange, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-error", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingErrorFmt() {
	// fmt error
	_, err := s.Client.PingError(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "fmt", ErrorCodeReturned: 1},
	)

	e := errors.FromError(err)
	assert.Equal(s.T(), errors.Internal, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-fmt", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingListSuccess() {
	// Success
	stream, err := s.Client.PingList(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "success"},
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	for i := 0; i < 5; i++ {
		_, err = stream.Recv()
		assert.NoError(s.T(), err, "reading stream should not fail")
	}

	_, err = stream.Recv()
	assert.Equal(s.T(), io.EOF, err, "stream should close with EOF")
}

func (s *errorInterceptorsSuite) TestPingListInternal() {
	// Internal
	stream, err := s.Client.PingList(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "internal", ErrorCodeReturned: 1},
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), "ping", e.GetComponent(), "Error component should have been set")
	assert.Equal(s.T(), "00100", e.Hex(), "Error code should be correct")
	assert.Equal(s.T(), "test-internal", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingListGRPC() {
	// gRPC
	stream, err := s.Client.PingList(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "test-error", ErrorCodeReturned: uint32(codes.OutOfRange)},
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), errors.OutOfRange, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-error", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingListFmt() {
	// gRPC
	stream, err := s.Client.PingList(
		s.SimpleCtx(),
		&testproto.PingRequest{Value: "fmt", ErrorCodeReturned: 1},
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), errors.Internal, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-fmt", e.GetMessage(), "Error message should be correct")
}

func (s *errorInterceptorsSuite) TestPingStreamSuccess() {
	stream, err := s.Client.PingStream(
		s.SimpleCtx(),
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	err = stream.Send(&testproto.PingRequest{Value: "success"})
	assert.NoError(s.T(), err, "writing to stream should not fail")

	_, err = stream.Recv()
	assert.NoError(s.T(), err, "reading stream should not fail")

	err = stream.CloseSend()
	assert.NoError(s.T(), err, "closing stream should not fail")
}

func (s *errorInterceptorsSuite) TestPingStreamInternal() {
	stream, err := s.Client.PingStream(
		s.SimpleCtx(),
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	err = stream.Send(&testproto.PingRequest{Value: "internal", ErrorCodeReturned: 1})
	assert.NoError(s.T(), err, "writing to stream should not fail")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), "ping", e.GetComponent(), "Error component should have been set")
	assert.Equal(s.T(), "00100", e.Hex(), "Error code should be correct")
	assert.Equal(s.T(), "test-internal", e.GetMessage(), "Error message should be correct")

	err = stream.CloseSend()
	assert.NoError(s.T(), err, "closing stream should not fail")
}

func (s *errorInterceptorsSuite) TestPingStreamGRPC() {
	stream, err := s.Client.PingStream(
		s.SimpleCtx(),
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	err = stream.Send(&testproto.PingRequest{Value: "test-error", ErrorCodeReturned: uint32(codes.InvalidArgument)})
	assert.NoError(s.T(), err, "writing to stream should not fail")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), errors.Data, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-error", e.GetMessage(), "Error message should be correct")

	err = stream.CloseSend()
	assert.NoError(s.T(), err, "closing stream should not fail")
}

func (s *errorInterceptorsSuite) TestPingStreamFmt() {
	stream, err := s.Client.PingStream(
		s.SimpleCtx(),
	)
	assert.NoError(s.T(), err, "should not fail on establishing the stream")

	err = stream.Send(&testproto.PingRequest{Value: "fmt", ErrorCodeReturned: 1})
	assert.NoError(s.T(), err, "writing to stream should not fail")

	_, err = stream.Recv()
	e := errors.FromError(err)
	assert.Equal(s.T(), errors.Internal, e.GetCode(), "Error code should be correct")
	assert.Equal(s.T(), "test-fmt", e.GetMessage(), "Error message should be correct")

	err = stream.CloseSend()
	assert.NoError(s.T(), err, "closing stream should not fail")
}

func newErrorInterceptorsSuite(t *testing.T) *errorInterceptorsSuite {
	return &errorInterceptorsSuite{
		InterceptorTestSuite: &grpctesting.InterceptorTestSuite{
			TestService: &errorPingService{&grpctesting.TestPingService{T: t}},
		},
	}
}

func TestInterceptors(t *testing.T) {
	b := newErrorInterceptorsSuite(t)
	b.InterceptorTestSuite.ClientOpts = []grpc.DialOption{
		grpc.WithUnaryInterceptor(UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(StreamClientInterceptor()),
	}
	b.InterceptorTestSuite.ServerOpts = []grpc.ServerOption{
		grpcmiddleware.WithUnaryServerChain(UnaryServerInterceptor()),
		grpcmiddleware.WithStreamServerChain(StreamServerInterceptor()),
	}
	suite.Run(t, b)
}
