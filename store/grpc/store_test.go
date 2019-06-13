package grpc

import (
	"context"
	"net"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/app/grpc/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type EnvelopeStoreTestSuite struct {
	server *grpc.Server
	conn   *grpc.ClientConn
	testutils.EnvelopeStoreTestSuite
}

func (s *EnvelopeStoreTestSuite) SetupTest() {
	s.server = grpc.NewServer()
	store.RegisterStoreServer(s.server, services.NewStoreService(mock.NewEnvelopeStore()))

	lis := bufconn.Listen(1024 * 1024)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	s.conn = conn
	s.Store = NewEnvelopeStore(store.NewStoreClient(conn))
}

func (s *EnvelopeStoreTestSuite) TearDownTest() {
	s.conn.Close()
	s.server.Stop()
}

func TestGRPC(t *testing.T) {
	s := new(EnvelopeStoreTestSuite)
	suite.Run(t, s)
}
