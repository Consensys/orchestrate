// +build unit

package controls

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
)

func TestCooldownControl_Execute(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	ctrl := NewCooldownControl()
	cooldown := 100 * time.Millisecond

	faucet1 := chainregistry.Faucet{
		UUID:       "001",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Cooldown:   cooldown,
	}
	faucet2 := chainregistry.Faucet{
		UUID:       "002",
		Amount:     big.NewInt(10),
		MaxBalance: big.NewInt(10),
		Cooldown:   cooldown,
	}

	t.Run("should skip first faucet on two consecutive requests", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*chainregistry.Request{
			newFaucetReq(candidates, chains[0], "", addresses[0]),
			newFaucetReq(candidates, chains[0], "", addresses[0]),
		}

		fcts := make([]chainregistry.Faucet, len(reqs))
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *chainregistry.Request) {
				err := ctrl.Control(ctx, req)
				assert.NoError(t, err)
				fct := electFirstFaucetCandidate(req.Candidates)
				err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
				ch <- true
				assert.NoError(t, err)
				fcts[idx] = fct
				wg.Done()
			}(idx, req)
			<-ch
		}

		wg.Wait()
		assert.NotEqual(t, fcts[0].UUID, fcts[1].UUID)
	})

	t.Run("should use same faucet in case cooldown time passed between requests", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*chainregistry.Request{
			newFaucetReq(candidates, chains[1], "", addresses[1]),
			newFaucetReq(candidates, chains[1], "", addresses[1]),
		}

		fcts := make([]chainregistry.Faucet, len(reqs))
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *chainregistry.Request) {
				err := ctrl.Control(ctx, req)
				assert.NoError(t, err)
				fct := electFirstFaucetCandidate(req.Candidates)
				err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
				ch <- true
				assert.NoError(t, err)
				fcts[idx] = fct
				wg.Done()
			}(idx, req)
			<-ch
			time.Sleep(cooldown)
		}

		wg.Wait()
		assert.Equal(t, fcts[0].UUID, fcts[0].UUID)
	})

	t.Run("should fail to request twice in less than cooldown requests", func(t *testing.T) {
		candidates := map[string]chainregistry.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*chainregistry.Request{
			newFaucetReq(candidates, chains[2], "", addresses[2]),
			newFaucetReq(candidates, chains[2], "", addresses[2]),
			newFaucetReq(candidates, chains[2], "", addresses[2]),
		}

		fcts := make([]chainregistry.Faucet, len(reqs))
		var eerr error
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *chainregistry.Request) {
				err := ctrl.Control(ctx, req)
				if err == nil {
					fct := electFirstFaucetCandidate(req.Candidates)
					err = ctrl.OnSelectedCandidate(ctx, &fct, req.Beneficiary)
					eerr = errors.CombineErrors(err)
					fcts[idx] = fct
				} else {
					eerr = err
				}
				ch <- true
				wg.Done()
			}(idx, req)
			<-ch
		}

		wg.Wait()
		assert.True(t, errors.IsFaucetWarning(eerr))
	})
}
