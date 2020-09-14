// +build unit

package controls

import (
	"math/big"
	"testing"
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
)

func TestMaxBalanceControl_Execute(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewMaxBalanceControl(client)

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   ethcommon.HexToAddress("0xcde"),
	}
	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(20),
		Creditor:   ethcommon.HexToAddress("0xeee"),
	}

	t.Run("should choose first candidate successfully", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(0), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
		assert.NoError(t, err)
		assert.Equal(t, fct.UUID, faucet1.UUID)
	})

	t.Run("should skip candidate first candidate because user exceeds max balance", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(10), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
		assert.NoError(t, err)
		assert.Equal(t, fct.UUID, faucet2.UUID)
	})

	t.Run("should fail to elect candidate because user exceeds max balances", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(20), nil)
		err := ctrl.Control(ctx, req)
		assert.True(t, errors.IsWarning(err))
	})

	t.Run("should fail when fetch balance fails", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(nil, errors.ConnectionError("cannot connect"))
		err := ctrl.Control(ctx, req)
		assert.NotNil(t, err)
	})
}
