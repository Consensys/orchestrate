// +build integration

package integrationtests

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"testing"
	"time"
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
}
