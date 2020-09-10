package chains

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/mock"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestRegisterChain_Execute(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)
	ethClient := mockethclient.NewMockClient(mockCtrl)

	ctx := context.Background()
	chainUUID := uuid.Must(uuid.NewV4()).String()
	chainID := big.NewInt(666)
	chainTip := big.NewInt(888)
	urls := []string{"http://geth1:8545", "http://geth2:8545"}

	registerChainUC := NewRegisterChain(chainAgent, ethClient)

	t.Run("should execute use case successfully and fetch chain tip", func(t *testing.T) {
		chain := &models.Chain{
			Name: "geth",
			URLs: urls,
		}

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().HeaderByNumber(gomock.Any(), urls[0], nil).Return(&ethtypes.Header{Number: chainTip}, nil)
		chainAgent.EXPECT().RegisterChain(ctx, chain).Return(nil)

		err := registerChainUC.Execute(ctx, chain)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully at defined starting block", func(t *testing.T) {
		startingBlock := uint64(666)
		chain := &models.Chain{
			UUID:                  chainUUID,
			Name:                  "geth",
			URLs:                  urls,
			ListenerStartingBlock: &startingBlock,
		}

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		chainAgent.EXPECT().RegisterChain(ctx, chain).Return(nil)

		err := registerChainUC.Execute(ctx, chain)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if Network fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		chain := &models.Chain{
			Name: "geth",
			URLs: urls,
		}

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(nil, expectedErr)

		err := registerChainUC.Execute(ctx, chain)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerChainComponent), err)
	})

	t.Run("should fail if different chainIDs are returned", func(t *testing.T) {
		chain := &models.Chain{
			Name: "geth",
			URLs: urls,
		}

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(big.NewInt(111), nil)

		err := registerChainUC.Execute(ctx, chain)

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if data agent fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		chain := &models.Chain{
			Name: "geth",
			URLs: urls,
		}

		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[1]).Return(chainID, nil)
		ethClient.EXPECT().Network(gomock.Any(), urls[0]).Return(chainID, nil)
		ethClient.EXPECT().HeaderByNumber(gomock.Any(), urls[0], nil).Return(&ethtypes.Header{Number: chainTip}, nil)
		chainAgent.EXPECT().RegisterChain(ctx, chain).Return(expectedErr)

		err := registerChainUC.Execute(ctx, chain)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerChainComponent), err)
	})
}
