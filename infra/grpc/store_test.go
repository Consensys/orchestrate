package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/grpc/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/testutils"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type TraceStoreTestSuite struct {
	server *grpc.Server
	conn   *grpc.ClientConn
	testutils.TraceStoreTestSuite
}

func (suite *TraceStoreTestSuite) SetupTest() {
	suite.server = grpc.NewServer()
	store.RegisterStoreServer(suite.server, services.NewStoreService(mock.NewTraceStore()))

	lis := bufconn.Listen(1024 * 1024)
	go func() {
		if err := suite.server.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithDialer(func(string, time.Duration) (net.Conn, error) { return lis.Dial() }),
		grpc.WithInsecure(),
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to dial bufnet: %v", err))
	}

	suite.conn = conn
	suite.Store = NewTraceStore(store.NewStoreClient(conn))
}

func (suite *TraceStoreTestSuite) TearDownTest() {
	suite.conn.Close()
	suite.server.Stop()
}

func TestGRPC(t *testing.T) {
	s := new(TraceStoreTestSuite)
	suite.Run(t, s)
}
