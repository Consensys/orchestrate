// +build unit

package controls

import (
	"context"
	"github.com/consensys/quorum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
)

func TestMaxBalanceControl_Execute(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewMaxBalanceControl(client)

	faucet1 := testutils.FakeFaucet()
	faucet1.Amount = "10"
	faucet1.MaxBalance = "20"
	faucet1.CreditorAccount = common.HexToAddress("0xcde").Hex()

	faucet2 := testutils.FakeFaucet()
	faucet2.Amount = "10"
	faucet2.MaxBalance = "20"
	faucet2.CreditorAccount = common.HexToAddress("0xeee").Hex()

	t.Run("should choose first candidate successfully", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(0), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, req.Candidates, 2)
	})

	t.Run("should skip first candidate because user exceeds max balance", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(10), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)

		assert.Len(t, req.Candidates, 2)
	})

	t.Run("should exclude all candidates where account exceeds max balance", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(20), nil)

		err := ctrl.Control(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, req.Candidates)
	})

	t.Run("should fail when fetch balance fails", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(nil, errors.ConnectionError("cannot connect"))
		err := ctrl.Control(ctx, req)
		assert.NotNil(t, err)
	})
}
