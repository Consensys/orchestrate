package usecases

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
)

func TestImportChain_Execute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockClient(mockCtrl)

	ctx := context.Background()
	chainID := big.NewInt(666)
	chainTip := big.NewInt(888)
	urls := []string{"http://geth1:8545", "http://geth2:8545"}
	importChainJSON := `{"name":"geth","urls":["http://geth1:8545", "http://geth2:8545"]}`

	registerChainUC := NewImportChain(chainAgent, ethClient)

	t.Run("should execute use case successfully and fetch chain tip", func(t *testing.T) {
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().HeaderByNumber(gomock.Any(), urls[0], nil).Return(&ethtypes.Header{Number: chainTip}, nil)
		chainAgent.EXPECT().RegisterChain(ctx, gomock.Any()).Return(nil)

		err := registerChainUC.Execute(ctx, importChainJSON)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully at defined starting block", func(t *testing.T) {
		chainJSON := `{"name":"geth","urls":["http://geth1:8545", "http://geth2:8545"], "listenerStartingBlock": "888"}`

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		chainAgent.EXPECT().RegisterChain(ctx, gomock.Any()).Return(nil)

		err := registerChainUC.Execute(ctx, chainJSON)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if Network fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(nil, expectedErr)

		err := registerChainUC.Execute(ctx, importChainJSON)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(importChainComponent), err)
	})

	t.Run("should fail with same error if data agent fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().HeaderByNumber(gomock.Any(), urls[0], nil).Return(&ethtypes.Header{Number: chainTip}, nil)
		chainAgent.EXPECT().RegisterChain(ctx, gomock.Any()).Return(expectedErr)

		err := registerChainUC.Execute(ctx, importChainJSON)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(importChainComponent), err)
	})
}
