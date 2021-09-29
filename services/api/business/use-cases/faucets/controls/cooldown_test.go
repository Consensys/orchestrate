// +build unit

package controls

import (
	"context"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
)

func TestCooldownControl_Execute(t *testing.T) {
	ctx := context.Background()

	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create CoolDown controlled credit
	ctrl := NewCooldownControl()
	cooldown := 100 * time.Millisecond

	faucet1 := testutils.FakeFaucet()
	faucet2 := testutils.FakeFaucet()

	t.Run("should skip first faucet on two consecutive requests", func(t *testing.T) {
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*entities.FaucetRequest{
			newFaucetReq(candidates, chains[0], "", addresses[0]),
			newFaucetReq(candidates, chains[0], "", addresses[0]),
		}

		fcts := make([]*entities.Faucet, len(reqs))
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *entities.FaucetRequest) {
				err := ctrl.Control(ctx, req)
				assert.NoError(t, err)
				fct := electFirstFaucetCandidate(req.Candidates)
				err = ctrl.OnSelectedCandidate(ctx, fct, req.Beneficiary)
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
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*entities.FaucetRequest{
			newFaucetReq(candidates, chains[1], "", addresses[1]),
			newFaucetReq(candidates, chains[1], "", addresses[1]),
		}

		fcts := make([]*entities.Faucet, len(reqs))
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *entities.FaucetRequest) {
				err := ctrl.Control(ctx, req)
				assert.NoError(t, err)
				fct := electFirstFaucetCandidate(req.Candidates)
				err = ctrl.OnSelectedCandidate(ctx, fct, req.Beneficiary)
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
		candidates := map[string]*entities.Faucet{
			faucet1.UUID: faucet1,
			faucet2.UUID: faucet2,
		}
		reqs := []*entities.FaucetRequest{
			newFaucetReq(candidates, chains[2], "", addresses[2]),
			newFaucetReq(candidates, chains[2], "", addresses[2]),
			newFaucetReq(candidates, chains[2], "", addresses[2]),
		}

		var eerr error
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)
		for idx, req := range reqs {
			wg.Add(1)
			go func(idx int, req *entities.FaucetRequest) {
				err := ctrl.Control(ctx, req)
				if err == nil {
					fct := electFirstFaucetCandidate(req.Candidates)
					err = ctrl.OnSelectedCandidate(ctx, fct, req.Beneficiary)
					eerr = errors.CombineErrors(err)
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
