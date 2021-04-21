// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/keymanager"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/zk-snarks"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/services/key-manager/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// keyManagerEthereumTestSuite is a test suite for Key Manager Ethereum controller
type keyManagerZKSTestSuite struct {
	suite.Suite
	baseURL string
	client  client.KeyManagerClient
	env     *IntegrationEnvironment
}

func (s *keyManagerZKSTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL, nil)
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *keyManagerZKSTestSuite) TestKeyManager_ZKS_CreateAndGet() {
	ctx := s.env.ctx

	accountRequest := testutils.FakeCreateZKSAccountRequest()
	account := &types.ZKSAccountResponse{}
	s.T().Run("should create a new account successfully", func(t *testing.T) {
		var err error
		account, err = s.client.ZKSCreateAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.NotEmpty(t, account.PublicKey)
		assert.Equal(t, entities.ZKSAlgorithmEDDSA, account.SigningAlgorithm)
		assert.Equal(t, entities.ZKSCurveBN254, account.Curve)
		assert.Equal(t, accountRequest.Namespace, account.Namespace)
	})

	s.T().Run("should get created account successfully", func(t *testing.T) {
		account2, err := s.client.ZKSGetAccount(ctx, account.PublicKey, account.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, account.PublicKey, account2.PublicKey)
	})
}

func (s *keyManagerZKSTestSuite) TestKeyManager_ZKS_Sign() {
	ctx := s.env.ctx
	accountRequest := testutils.FakeCreateZKSAccountRequest()
	account, err := s.client.ZKSCreateAccount(ctx, accountRequest)
	assert.NoError(s.T(), err)

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		signRequest := &keymanager.SignPayloadRequest{
			Data:      "0xda",
			Namespace: account.Namespace,
		}
		signature, err := s.client.ZKSSign(ctx, account.PublicKey, signRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := &keymanager.SignPayloadRequest{
			Data: "",
		}

		_, err := s.client.ZKSSign(ctx, account.PublicKey, signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})
}

func (s *keyManagerZKSTestSuite) TestKeyManager_ZKS_VerifySignature() {
	ctx := s.env.ctx

	accountRequest := testutils.FakeCreateZKSAccountRequest()
	account, err := s.client.ZKSCreateAccount(ctx, accountRequest)
	require.NoError(s.T(), err)

	signRequest := &keymanager.SignPayloadRequest{
		Data:      "0xda",
		Namespace: account.Namespace,
	}

	signature, err := s.client.ZKSSign(ctx, account.PublicKey, signRequest)
	require.NoError(s.T(), err)

	s.T().Run("should verify signature successfully", func(t *testing.T) {
		err = s.client.ZKSVerifySignature(ctx, &types.VerifyPayloadRequest{
			Data:      signRequest.Data,
			Signature: signature,
			PublicKey: account.PublicKey,
		})

		assert.NoError(t, err)
	})

	s.T().Run("should fail with 400 if payload is invalid", func(t *testing.T) {
		err := s.client.ZKSVerifySignature(ctx, &types.VerifyPayloadRequest{
			Data:      signRequest.Data,
			Signature: "",
			PublicKey: account.PublicKey,
		})

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 422 if signature is invalid", func(t *testing.T) {
		err := s.client.ZKSVerifySignature(ctx, &types.VerifyPayloadRequest{
			Data:      signRequest.Data,
			Signature: "0xfeee",
			PublicKey: account.PublicKey,
		})

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with 422 if signature is invalid for the given address", func(t *testing.T) {
		err := s.client.ZKSVerifySignature(ctx, &types.VerifyPayloadRequest{
			Data:      signRequest.Data,
			Signature: signature,
			PublicKey: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		})

		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
