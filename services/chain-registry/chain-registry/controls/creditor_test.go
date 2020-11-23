// +build unit

package controls

import (
	"context"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
)

func TestCreditorControl_SuccessfulCandidate(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	client := mock.NewMockChainStateReader(mockCtrl)
	ctrl := NewCreditorControl(client)

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[0],
	}

	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[1],
	}

	t.Run("should choose first candidate successfully", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(10), nil)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[1], nil).Return(big.NewInt(10), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
		assert.NoError(t, err)
		assert.Equal(t, fct.UUID, faucet1.UUID)
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

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[0],
	}
	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[1],
	}

	t.Run("should skip candidate when beneficiary is same as creditor", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[0])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[1], nil).Return(big.NewInt(10), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
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

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[0],
	}
	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[1],
	}

	t.Run("should skip candidate when creditor does not have enough balance", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(0), nil)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[1], nil).Return(big.NewInt(10), nil)
		err := ctrl.Control(ctx, req)
		assert.NoError(t, err)
		fct := electFirstFaucetCandidate(req.Candidates)
		err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
		assert.NoError(t, err)
		assert.Equal(t, fct.UUID, faucet2.UUID)
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

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[0],
	}
	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Creditor:   addresses[1],
	}

	t.Run("should fail with no candidate creditor with enough balance", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).Return(big.NewInt(0), nil)
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[1], nil).Return(big.NewInt(0), nil)
		err := ctrl.Control(ctx, req)
		assert.True(t, errors.IsWarning(err))
	})

	t.Run("should fail when fetch balance fails", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
		}
		req := newFaucetReq(candidates, chains[0], chainURLs[0], addresses[2])
		client.EXPECT().BalanceAt(gomock.Any(), gomock.Any(), addresses[0], nil).
			Return(nil, errors.ConnectionError("cannot connect"))
		err := ctrl.Control(ctx, req)
		assert.NotNil(t, err)
	})
}
