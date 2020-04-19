// +build unit

package grpcmetrics

import (
	"testing"

	"github.com/golang/mock/gomock"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpctesting "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mockmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
	"google.golang.org/grpc"
)

type errorInterceptorsSuite struct {
	suite.Suite
	InterceptorTestSuite *grpctesting.InterceptorTestSuite
	ctrlr                *gomock.Controller
	registry             *mockmetrics.MockGRPCServer
}

func (s *errorInterceptorsSuite) SetupTest() {
	s.ctrlr = gomock.NewController(s.T())
	s.registry = mockmetrics.NewMockGRPCServer(s.ctrlr)
	s.InterceptorTestSuite = &grpctesting.InterceptorTestSuite{
		TestService: &grpctesting.TestPingService{T: s.T()},
		ServerOpts: []grpc.ServerOption{
			grpcmiddleware.WithUnaryServerChain(UnaryServerInterceptor(s.registry)),
			grpcmiddleware.WithStreamServerChain(StreamServerInterceptor(s.registry)),
		},
	}
	s.InterceptorTestSuite.SetT(s.T())
	s.InterceptorTestSuite.SetupSuite()
}

func (s *errorInterceptorsSuite) TearDownTest() {
	s.InterceptorTestSuite.TearDownSuite()
	s.ctrlr.Finish()
}
func (s *errorInterceptorsSuite) TestPing() {
	startedCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StartedCounter().Return(startedCounter)
	startedCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "unary", "method", "Ping").
		Return(startedCounter)
	startedCounter.EXPECT().Add(float64(1))

	handledCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().HandledCounter().Return(handledCounter)
	handledCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "unary", "method", "Ping", "code", "OK").
		Return(handledCounter)
	handledCounter.EXPECT().Add(float64(1))

	handledHistogram := mockmetrics.NewMockHistogram(s.ctrlr)
	s.registry.EXPECT().HandledDurationHistogram().Return(handledHistogram)
	handledHistogram.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "unary", "method", "Ping", "code", "OK").
		Return(handledHistogram)
	handledHistogram.EXPECT().Observe(gomock.Any())

	_, err := s.InterceptorTestSuite.Client.Ping(
		s.InterceptorTestSuite.SimpleCtx(),
		// context.Background(),
		&testproto.PingRequest{Value: "test"},
	)
	assert.NoError(s.T(), err)
}

func (s *errorInterceptorsSuite) TestPingList() {
	startedCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StartedCounter().Return(startedCounter)
	startedCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "server_stream", "method", "PingList").
		Return(startedCounter)
	startedCounter.EXPECT().Add(float64(1))

	handledCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().HandledCounter().Return(handledCounter)
	handledCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "server_stream", "method", "PingList", "code", "OK").
		Return(handledCounter)
	handledCounter.EXPECT().Add(float64(1))

	handledHistogram := mockmetrics.NewMockHistogram(s.ctrlr)
	s.registry.EXPECT().HandledDurationHistogram().Return(handledHistogram)
	handledHistogram.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "server_stream", "method", "PingList", "code", "OK").
		Return(handledHistogram)
	handledHistogram.EXPECT().Observe(gomock.Any())

	receivedMsgCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StreamMsgReceivedCounter().Return(receivedMsgCounter)
	receivedMsgCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "server_stream", "method", "PingList").
		Return(receivedMsgCounter)
	receivedMsgCounter.EXPECT().Add(float64(1))

	sentMsgCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StreamMsgSentCounter().Return(sentMsgCounter)
	sentMsgCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "server_stream", "method", "PingList").
		Return(sentMsgCounter)
	sentMsgCounter.EXPECT().Add(float64(1)).Times(100)

	_, err := s.InterceptorTestSuite.Client.PingList(
		s.InterceptorTestSuite.SimpleCtx(),
		&testproto.PingRequest{Value: "test"},
	)
	assert.NoError(s.T(), err)
}

func (s *errorInterceptorsSuite) TestPingStream() {
	startedCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StartedCounter().Return(startedCounter)
	startedCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "bidi_stream", "method", "PingStream").
		Return(startedCounter)
	startedCounter.EXPECT().Add(float64(1))

	handledCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().HandledCounter().Return(handledCounter)
	handledCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "bidi_stream", "method", "PingStream", "code", "OK").
		Return(handledCounter)
	handledCounter.EXPECT().Add(float64(1))

	handledHistogram := mockmetrics.NewMockHistogram(s.ctrlr)
	s.registry.EXPECT().HandledDurationHistogram().Return(handledHistogram)
	handledHistogram.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "bidi_stream", "method", "PingStream", "code", "OK").
		Return(handledHistogram)
	handledHistogram.EXPECT().Observe(gomock.Any())

	receivedMsgCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StreamMsgReceivedCounter().Return(receivedMsgCounter)
	receivedMsgCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "bidi_stream", "method", "PingStream").
		Return(receivedMsgCounter)
	receivedMsgCounter.EXPECT().Add(float64(1)).Times(5)

	sentMsgCounter := mockmetrics.NewMockCounter(s.ctrlr)
	s.registry.EXPECT().StreamMsgSentCounter().Return(sentMsgCounter)
	sentMsgCounter.EXPECT().
		With("tenant_id", "_", "service", "mwitkow.testproto.TestService", "type", "bidi_stream", "method", "PingStream").
		Return(sentMsgCounter)
	sentMsgCounter.EXPECT().Add(float64(1)).Times(5)

	stream, err := s.InterceptorTestSuite.Client.PingStream(
		s.InterceptorTestSuite.SimpleCtx(),
	)
	assert.NoError(s.T(), err)

	for i := 0; i < 5; i++ {
		err = stream.Send(&testproto.PingRequest{Value: "test"})
		assert.NoError(s.T(), err, "writing to stream should not fail")
	}

	err = stream.CloseSend()
	assert.NoError(s.T(), err, "closing stream should not fail")

	for i := 0; i < 5; i++ {
		_, err = stream.Recv()
		assert.NoError(s.T(), err, "reading from stream should not fail")
	}

	_, _ = stream.Recv()
}

func TestInterceptors(t *testing.T) {
	suite.Run(t, &errorInterceptorsSuite{})
}
