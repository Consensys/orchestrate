// +build unit

package controls

import (
	"context"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"math/big"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreditorControl_SuccessfulCandidate(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewCreditorControl(client)

	t.Run("should choose first candidate successfully", func(t *testing.T) {
		faucet1 := testutils.FakeFaucet()
		faucet1.CreditorAccount = addresses[0]
		faucet2 := testutils.FakeFaucet()
		faucet2.CreditorAccount = addresses[1]

		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])

		gomock.InOrder(
			client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(10000), nil),
			client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(10000), nil),
		)

		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)

		assert.Len(t, req.Candidates, 2)
	})
}

func TestCreditorControl_SkipBeneficiaryAsCandidate(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewCreditorControl(client)

	t.Run("should skip candidate when beneficiary is same as creditor", func(t *testing.T) {
		faucet1 := testutils.FakeFaucet()
		faucet1.CreditorAccount = addresses[0]
		faucet2 := testutils.FakeFaucet()
		faucet2.CreditorAccount = addresses[1]

		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], faucet1.CreditorAccount)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(1000000), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, fct, req.Beneficiary)
		assert.NoError(t, err)
		assert.Equal(t, fct.UUID, faucet2.UUID)
	})
}

func TestCreditorControl_SkipCandidateWithNotFounds(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewCreditorControl(client)

	faucet1 := testutils.FakeFaucet()
	faucet2 := testutils.FakeFaucet()

	t.Run("should skip candidate when creditor does not have enough balance", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(0), nil)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(1000000), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, req.Candidates, 1)
	})
}

func TestCreditorControl_Failure(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewCreditorControl(client)

	faucet1 := testutils.FakeFaucet()
	faucet2 := testutils.FakeFaucet()

	t.Run("should remove all candidates without enough balance", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(0), nil)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(big.NewInt(0), nil)

		err := ctrl.Control(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, req.Candidates)
	})

	t.Run("should fail when fetch balance fails", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), gomock.Any(), nil).
			Return(nil, errors.ConnectionError("cannot connect"))
		err := ctrl.Control(ctx, req)
		assert.NotNil(t, err)
	})
}
