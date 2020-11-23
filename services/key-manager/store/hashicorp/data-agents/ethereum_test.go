// +build unit

package dataagents

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/secretstore/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ethereumDATestSuite struct {
	suite.Suite
	dataAgent       *HashicorpEthereum
	mockSecretStore *mocks.MockSecretStore
	router          *mux.Router
}

func TestEthereumDataAgent(t *testing.T) {
	s := new(ethereumDATestSuite)
	suite.Run(t, s)
}

func (s *ethereumDATestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockSecretStore = mocks.NewMockSecretStore(ctrl)

	s.dataAgent = NewHashicorpEthereum(s.mockSecretStore)
}

func (s *ethereumDATestSuite) TestEthereumDataAgent_Insert() {
	ctx := context.Background()
	address := "0xaddress"
	privKey := "privKey"
	namespace := "namespace"

	s.T().Run("should insert private key successfully without namespace", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Store(ctx, address, privKey).Return(nil)

		err := s.dataAgent.Insert(ctx, address, privKey, "")
		assert.NoError(t, err)
	})

	s.T().Run("should insert private key successfully with namespace", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Store(ctx, "namespace0xaddress", privKey).Return(nil)

		err := s.dataAgent.Insert(ctx, address, privKey, namespace)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with HashicorpVaultConnectionError if Store fails", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Store(ctx, gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

		err := s.dataAgent.Insert(ctx, address, privKey, namespace)
		assert.Equal(t, errors.HashicorpVaultConnectionError("failed to store privateKey in Hashicorp Vault").ExtendComponent(ethereumDAComponent), err)
	})
}

func (s *ethereumDATestSuite) TestEthereumDataAgent_FindOne() {
	ctx := context.Background()
	address := "0xaddress"
	privKey := "privKey"
	namespace := "namespace"

	s.T().Run("should insert private key successfully without namespace", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Load(ctx, address).Return(privKey, true, nil)

		privKeyResponse, err := s.dataAgent.FindOne(ctx, address, "")

		assert.NoError(t, err)
		assert.Equal(t, privKey, privKeyResponse)
	})

	s.T().Run("should insert private key successfully with namespace", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Load(ctx, "namespace0xaddress").Return(privKey, true, nil)

		privKeyResponse, err := s.dataAgent.FindOne(ctx, address, namespace)

		assert.NoError(t, err)
		assert.Equal(t, privKey, privKeyResponse)
	})

	s.T().Run("should fail with HashicorpVaultConnectionError if Load fails", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Load(ctx, gomock.Any()).Return("", false, fmt.Errorf("error"))

		privKeyResponse, err := s.dataAgent.FindOne(ctx, address, namespace)

		assert.Empty(t, privKeyResponse)
		assert.Equal(t, errors.HashicorpVaultConnectionError("failed to load privateKey from Hashicorp Vault").ExtendComponent(ethereumDAComponent), err)
	})

	s.T().Run("should fail with NotFoundError if Load succeeds but nothing is returned", func(t *testing.T) {
		s.mockSecretStore.EXPECT().Load(ctx, gomock.Any()).Return("", false, nil)

		privKeyResponse, err := s.dataAgent.FindOne(ctx, address, namespace)

		assert.Empty(t, privKeyResponse)
		assert.Equal(t, errors.NotFoundError("account does not exist").ExtendComponent(ethereumDAComponent), err)
	})
}
