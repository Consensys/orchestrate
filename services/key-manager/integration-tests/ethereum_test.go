// +build integration

package integrationtests

import (
	"github.com/consensys/quorum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"
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

		account, err := s.client.ETHCreateAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.True(t, common.IsHexAddress(account.Address))
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Import() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = ""

		_, err := s.client.ETHImportAccount(ctx, accountRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should create a new account successfully", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		account, err := s.client.ETHImportAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Sign() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := &keymanager.PayloadRequest{
			Data: "",
		}

		_, err := s.client.ETHSign(ctx, "0xaddress", signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"

		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		account, err := s.client.ETHImportAccount(ctx, accountRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, account.Address)

		signRequest := &keymanager.PayloadRequest{
			Data:      "my data to sign",
			Namespace: "_",
		}
		signature, err := s.client.ETHSign(ctx, expectedAddress, signRequest)
		assert.NoError(t, err)
		assert.Equal(t, "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00", signature)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_SignTransaction() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := testutils.FakeSignETHTransactionRequest()
		signRequest.ChainID = ""

		_, err := s.client.ETHSignTransaction(ctx, "0xaddress", signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign transaction successfully", func(t *testing.T) {
		expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"

		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		account, err := s.client.ETHImportAccount(ctx, accountRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, account.Address)

		signRequest := testutils.FakeSignETHTransactionRequest()
		signature, err := s.client.ETHSignTransaction(ctx, expectedAddress, signRequest)
		assert.NoError(t, err)
		assert.Equal(t, "0x3dcedc00acdd28aab04c2b352608fc6a3cbb17a82935c9168f434ee6d85ddbdd6c75f3299b37977796c019825c9ef49626fd805daa46efc495c5abb2e836446a01", signature)
	})
}
