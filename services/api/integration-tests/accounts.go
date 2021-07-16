// +build integration

package integrationtests

import (
	"fmt"
	"testing"
	"time"

	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type accountsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *accountsTestSuite) TestCreateAccounts() {
	ctx := s.env.ctx
	chain := testutils.FakeChain()
	chain.URLs = []string{s.env.blockchainNodeURL}

	s.T().Run("should create account successfully by querying key-manager API", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		resp, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), resp.Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.TenantID, "_")
	})

	s.T().Run("should fail to create account with same alias", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		_, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)
		log.WithoutContext().Errorf("%v", err)
		assert.True(s.T(), errors.IsAlreadyExistsError(err))
	})

	s.T().Run("should create account successfully and trigger funding transaction", func(t *testing.T) {
		chainWithFaucet, err := s.client.RegisterChain(s.env.ctx, &api.RegisterChainRequest{
			Name: "ganache-with-faucet-accounts",
			URLs: []string{s.env.blockchainNodeURL},
			Listener: api.RegisterListenerRequest{
				FromBlock:         "latest",
				ExternalTxEnabled: false,
			},
		})
		require.NoError(s.T(), err)

		privKey, address, err := createNewKey()
		require.NoError(s.T(), err, fmt.Sprintf("invalid private key %s", privKey))
		accountFaucet := testutils.FakeAccount()
		accountFaucet.Alias = "MyFaucetCreditor-accounts_" + utils.RandString(5)
		accountFaucet.Address = address

		req := testutils.FakeImportAccountRequest()
		req.PrivateKey = privKey
		req.Alias = accountFaucet.Alias
		ethAccRes, err := s.client.ImportAccount(s.env.ctx, req)
		require.NoError(s.T(), err)

		faucetRequest := testutils.FakeRegisterFaucetRequest()
		faucetRequest.Name = "faucet-integration-tests"
		faucetRequest.ChainRule = chainWithFaucet.UUID
		faucetRequest.CreditorAccount = accountFaucet.Address
		faucet, err := s.client.RegisterFaucet(s.env.ctx, faucetRequest)
		require.NoError(s.T(), err)

		accountRequest := testutils.FakeCreateAccountRequest()
		accountRequest.Chain = chainWithFaucet.Name

		require.NoError(s.T(), err)

		assert.Equal(s.T(), ethAccRes.TenantID, "_")

		err = s.client.DeleteChain(ctx, chainWithFaucet.UUID)
		assert.NoError(s.T(), err)
		err = s.client.DeleteFaucet(ctx, faucet.UUID)
		assert.NoError(s.T(), err)
	})

	s.T().Run("should fail to create account if postgres is down", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(s.T(), err)

		_, err = s.client.CreateAccount(ctx, txRequest)
		assert.Error(s.T(), err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)
	})
}

func (s *accountsTestSuite) TestImport() {
	ctx := s.env.ctx

	s.T().Run("should import account successfully by querying key-manager API", func(t *testing.T) {
		privKey, _, err := createNewKey()
		require.NoError(s.T(), err, fmt.Sprintf("invalid private key %s", privKey))
		txRequest := testutils.FakeImportAccountRequest()
		txRequest.PrivateKey = privKey

		resp, err := s.client.ImportAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), resp.Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.TenantID, "_")
	})

	s.T().Run("should fail to import same account twice", func(t *testing.T) {
		privKey, _, err := createNewKey()
		require.NoError(s.T(), err, fmt.Sprintf("invalid private key %s", privKey))
		txRequest := testutils.FakeImportAccountRequest()
		txRequest.PrivateKey = privKey

		_, err = s.client.ImportAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		privKey, _, err = createNewKey()
		require.NoError(s.T(), err)
		txRequest.PrivateKey = privKey
		_, err = s.client.ImportAccount(ctx, txRequest)
		require.Error(s.T(), err)
		assert.True(s.T(), errors.IsAlreadyExistsError(err))
	})
}

func (s *accountsTestSuite) TestSearch() {
	ctx := s.env.ctx

	s.T().Run("should create account and search for it by alias successfully", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		resp, err := s.client.SearchAccounts(ctx, &entities.AccountFilters{
			Aliases: []string{txRequest.Alias},
		})
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Len(s.T(), resp, 1)
		assert.Equal(s.T(), resp[0].Address, ethAccRes.Address)
		assert.Equal(s.T(), resp[0].PublicKey, ethAccRes.PublicKey)
		assert.Equal(s.T(), resp[0].Alias, txRequest.Alias)
		assert.Equal(s.T(), resp[0].TenantID, "_")
	})
}

func (s *accountsTestSuite) TestGetOne() {
	ctx := s.env.ctx

	s.T().Run("should create account and get it by address successfully", func(t *testing.T) {
		defer gock.Off()
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		resp, err := s.client.GetAccount(ctx, ethAccRes.Address)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), resp.Address, ethAccRes.Address)
		assert.Equal(s.T(), resp.PublicKey, ethAccRes.PublicKey)
		assert.Equal(s.T(), resp.Alias, txRequest.Alias)
		assert.Equal(s.T(), resp.TenantID, "_")
	})
}

func (s *accountsTestSuite) TestUpdate() {
	ctx := s.env.ctx

	s.T().Run("should create account and update it successfully", func(t *testing.T) {
		txRequest := testutils.FakeCreateAccountRequest()

		ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		txRequest2 := testutils.FakeUpdateAccountRequest()
		resp, err := s.client.UpdateAccount(ctx, ethAccRes.Address, txRequest2)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), resp.Alias, txRequest2.Alias)
		assert.Equal(s.T(), resp.Attributes, txRequest2.Attributes)
		assert.Equal(s.T(), resp.TenantID, "_")
	})
}

func (s *accountsTestSuite) TestSignPayload() {
	ctx := s.env.ctx
	txRequest := testutils.FakeCreateAccountRequest()
	ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
	require.NoError(s.T(), err)

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		address := ethAccRes.Address
		payload := hexutil.Encode([]byte("my data to sign"))

		signedPayload, err := s.client.SignPayload(ctx, address, &api.SignPayloadRequest{
			Data: payload,
		})
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedPayload)
	})
}

func (s *accountsTestSuite) TestVerifySignature() {
	ctx := s.env.ctx

	s.T().Run("should verify signature successfully", func(t *testing.T) {
		verifyRequest := qkm.FakeVerifyPayloadRequest()
		err := s.client.VerifySignature(ctx, verifyRequest)
		assert.NoError(s.T(), err)
	})
}

func (s *accountsTestSuite) TestSignTypedData() {
	ctx := s.env.ctx

	txRequest := testutils.FakeCreateAccountRequest()
	ethAccRes, err := s.client.CreateAccount(ctx, txRequest)
	require.NoError(s.T(), err)

	typedDataRequest := qkm.FakeSignTypedDataRequest()
	var signature string

	s.T().Run("should sign typed data successfully", func(t *testing.T) {
		signature, err = s.client.SignTypedData(ctx, ethAccRes.Address, &api.SignTypedDataRequest{
			DomainSeparator: typedDataRequest.DomainSeparator,
			Types:           typedDataRequest.Types,
			Message:         typedDataRequest.Message,
			MessageType:     typedDataRequest.MessageType,
		})

		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)
	})

	s.T().Run("should verify typed data signature successfully", func(t *testing.T) {
		err := s.client.VerifyTypedDataSignature(ctx, &qkmtypes.VerifyTypedDataRequest{
			TypedData: *typedDataRequest,
			Signature: hexutil.MustDecode(signature),
			Address:   common.HexToAddress(ethAccRes.Address),
		})
		assert.NoError(s.T(), err)
	})
}

func createNewKey() (string, string, error) {
	faucetKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privKey := hexutil.Encode(faucetKey.D.Bytes())
	address := crypto.PubkeyToAddress(faucetKey.PublicKey).String()
	return privKey[2:], address, nil
}
