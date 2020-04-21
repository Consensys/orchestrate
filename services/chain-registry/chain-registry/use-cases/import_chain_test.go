package usecases

import (
	"context"
	"fmt"
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

func TestImportChain_FetchHead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockChainLedgerReader(mockCtrl)

	importChainUC := NewImportChain(chainAgent, ethClient)
	uuid := genuuid.NewV4().String()
	importChainJSON := fmt.Sprintf(`{"uuid":"%s", "name":"geth","urls":["http://geth:8545"]}`, uuid)

	ethClient.EXPECT().HeaderByNumber(gomock.Any(), gomock.Eq("http://geth:8545"), nil).
		Return(&ethtypes.Header{Number: big.NewInt(666)}, nil).Times(1)

	expectedChain := &models.Chain{
		UUID:                  uuid,
		Name:                  "geth",
		URLs:                  []string{"http://geth:8545"},
		ListenerStartingBlock: &(&struct{ x uint64 }{666}).x,
	}
	expectedChain.SetDefault()
	chainAgent.EXPECT().RegisterChain(gomock.Any(), gomock.Eq(expectedChain))

	err := importChainUC.Execute(context.Background(), importChainJSON)
	assert.Nil(t, err)
}

func TestImportChain_NotFetchHead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockChainLedgerReader(mockCtrl)

	importChainUC := NewImportChain(chainAgent, ethClient)
	uuid := genuuid.NewV4().String()
	importChainJSON := fmt.Sprintf(`{"uuid":"%s", "name":"geth","urls":["http://geth:8545"],"listenerStartingBlock":"666"}`, uuid)

	expectedChain := &models.Chain{
		UUID:                  uuid,
		Name:                  "geth",
		URLs:                  []string{"http://geth:8545"},
		ListenerStartingBlock: &(&struct{ x uint64 }{666}).x,
	}
	expectedChain.SetDefault()
	chainAgent.EXPECT().RegisterChain(gomock.Any(), gomock.Eq(expectedChain))

	err := importChainUC.Execute(context.Background(), importChainJSON)
	assert.Nil(t, err)
}
