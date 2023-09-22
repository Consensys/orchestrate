//go:build integration
// +build integration

package integrationtests

import (
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

		//FIXME CUSTOM HEADER

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
}
