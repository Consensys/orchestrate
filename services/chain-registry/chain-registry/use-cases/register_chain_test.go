package usecases

import (
	"context"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mockethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestRegisterChain_FetchHead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockChainLedgerReader(mockCtrl)

	registerChainUC := NewRegisterChain(chainAgent, ethClient)
	uuid := genuuid.NewV4().String()

	ethClient.EXPECT().HeaderByNumber(gomock.Any(), gomock.Eq("http://geth:8545"), nil).
		Return(&ethtypes.Header{Number: big.NewInt(666)}, nil).Times(1)

	chain := &models.Chain{
		UUID: uuid,
		Name: "geth",
		URLs: []string{"http://geth:8545"},
	}

	expectedChain := *chain
	expectedChain.ListenerStartingBlock = &(&struct{ x uint64 }{666}).x
	expectedChain.SetDefault()

	chainAgent.EXPECT().RegisterChain(gomock.Any(), gomock.Eq(&expectedChain))

	err := registerChainUC.Execute(context.Background(), chain)
	assert.Nil(t, err)
}

func TestRegisterChain_NotFetchHead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockChainLedgerReader(mockCtrl)

	registerChainUC := NewRegisterChain(chainAgent, ethClient)
	uuid := genuuid.NewV4().String()

	chain := &models.Chain{
		UUID:                  uuid,
		Name:                  "geth",
		URLs:                  []string{"http://geth:8545"},
		ListenerStartingBlock: &(&struct{ x uint64 }{666}).x,
	}

	expectedChain := *chain
	expectedChain.SetDefault()

	chainAgent.EXPECT().RegisterChain(gomock.Any(), gomock.Eq(&expectedChain))

	err := registerChainUC.Execute(context.Background(), chain)
	assert.Nil(t, err)
}
