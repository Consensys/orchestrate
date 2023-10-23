//go:build integration
// +build integration

package integrationtests

import (
	http2 "net/http"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type proxyTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *proxyTestSuite) TestProxy() {
	ctx := s.env.ctx

	s.T().Run("should register chain and create proxy to the node", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.FromBlock = "latest"
		req.URLs = []string{s.env.blockchainNodeURL}

		chain, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		err = backoff.RetryNotify(
			func() error {
				_, der := ethclient.GlobalClient().Network(ctx, utils.GetProxyURL(s.env.baseURL, chain.UUID))
				return der
			},
			backoff.WithMaxRetries(backoff.NewConstantBackOff(2*time.Second), 5),
			func(_ error, _ time.Duration) {},
		)

		require.NoError(t, err)
	})

	s.T().Run("should register chain connected to a protected rpc npc and create proxy to the node", func(t *testing.T) {
		defer gock.OffAll()
		myFakeProtectedNode := "http://fakeNode.com"
		expectedHeaders := map[string]string{
			"x-api-key": "my-key",
		}
		req := testdata.FakeRegisterChainRequest()
		req.URL = myFakeProtectedNode
		req.Headers = expectedHeaders

		gock.New(myFakeProtectedNode).
			Post("").
			MatchHeader("x-api-key", "my-key").
			Reply(http2.StatusOK).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x6f\"}")

		chain, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)
		require.Equal(t, uint64(0x6f), chain.ChainID)

		// Following part it is not testable till we remove traefik logic for proxying traffic
		//
		// gock.New(myFakeProtectedNode).
		// 	Post("").
		// 	MatchHeader("x-api-key", "my-key").
		// 	Reply(http2.StatusOK).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x6f\"}")
		//
		// ec, err := ethclient.NewClient(fmt.Sprintf("%s/proxy/chains/%s", s.env.baseURL, chain.UUID))
		// require.NoError(t, err)
		// v, err := ec.ChainID(ctx)
		// require.NoError(t, err)
		// require.Equal(t, uint64(0x6f), v.Uint64())
	})
}
