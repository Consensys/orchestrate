// +build integration

package integrationtests

import (
	"github.com/consensys/quorum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	"testing"
)

// keyManagerEthereumTestSuite is a test suite for Key Manager Ethereum controller
type keyManagerEthereumTestSuite struct {
	suite.Suite
	baseURL string
	client  client.KeyManagerClient
	env     *IntegrationEnvironment
}

func (s *keyManagerEthereumTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL, nil)
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Create() {
	ctx := s.env.ctx

	s.T().Run("should create a new account successfully", func(t *testing.T) {
		accountRequest := testutils.FakeCreateETHAccountRequest()

		account, err := s.client.CreateETHAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.True(t, common.IsHexAddress(account.Address))
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Import() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = ""

		_, err := s.client.ImportETHAccount(ctx, accountRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should create a new account successfully", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		account, err := s.client.ImportETHAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})
}
