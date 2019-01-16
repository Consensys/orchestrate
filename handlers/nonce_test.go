package handlers

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/infra"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

const chainNonce uint64 = 42

func getChainNonce(withError bool) GetNonceFunc {
	if withError == true {
		return func(chainID *big.Int, a *common.Address) (uint64, error) {
			return 0, errors.New("error")
		}
	}

	return func(chainID *big.Int, a *common.Address) (uint64, error) {
		return chainNonce, nil
	}
}

func makeNonceContext() *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	return ctx
}

type MockNonceManager struct {
	mux   *sync.Mutex
	nonce uint64
}

func (nm *MockNonceManager) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return nm.nonce, true, nil
}

func (nm *MockNonceManager) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	nm.nonce = newNonce
	return nil
}

func (nm *MockNonceManager) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	nm.mux.Lock()
	return "random", nil
}

func (nm *MockNonceManager) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	nm.mux.Unlock()
	return nil
}

func TestNonceHandler(t *testing.T) {
	nm := MockNonceManager{
		mux:   &sync.Mutex{},
		nonce: 0,
	}
	nonceH := NonceHandler(&nm, getChainNonce(false))

	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeNonceContext()
		go func(ctx *types.Context) {
			defer wg.Done()
			nonceH(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("NonceHandler: expected %v outs but got %v", rounds, len(outs))
	}

	var n uint64
	nonceSet := make(map[uint64]bool)
	for ctx := range outs {
		n = ctx.T.Tx().Nonce()
		nonceSet[n] = true
	}
	for nonce := uint64(0); nonce < uint64(rounds); nonce++ {
		if nonceSet[nonce] == false {
			t.Errorf("NonceHandler: nonce %v is missing but should has been given", nonce)
		}
	}
}

type MockNonceManagerFromChain struct {
	mux   *sync.Mutex
	nonce uint64
}

func (nm *MockNonceManagerFromChain) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return nm.nonce, false, nil
}

func (nm *MockNonceManagerFromChain) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	nm.nonce = newNonce
	return nil
}

func (nm *MockNonceManagerFromChain) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	nm.mux.Lock()
	return "random", nil
}

func (nm *MockNonceManagerFromChain) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	nm.mux.Unlock()
	return nil
}

func TestNonceHandlerFromChain(t *testing.T) {
	nm := MockNonceManagerFromChain{
		mux:   &sync.Mutex{},
		nonce: 0,
	}
	nonceH := NonceHandler(&nm, getChainNonce(false))
	ctx := makeNonceContext()
	nonceH(ctx)
	if ctx.T.Tx().Nonce() != chainNonce {
		t.Errorf("NonceHandler: nonce should have come from getNonceChain")
	}
}

type MockNonceManagerGetNonceError struct{}

func (nm *MockNonceManagerGetNonceError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, true, errors.New("error")
}
func (nm *MockNonceManagerGetNonceError) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerGetNonceError) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerGetNonceError) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerUpdateCacheNonceError struct{}

func (nm *MockNonceManagerUpdateCacheNonceError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, true, nil
}
func (nm *MockNonceManagerUpdateCacheNonceError) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return errors.New("error")
}
func (nm *MockNonceManagerUpdateCacheNonceError) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerUpdateCacheNonceError) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerGetLockError struct{}

func (nm *MockNonceManagerGetLockError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, true, nil
}
func (nm *MockNonceManagerGetLockError) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerGetLockError) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", errors.New("error")
}
func (nm *MockNonceManagerGetLockError) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerGetChainNonceError struct{}

func (nm *MockNonceManagerGetChainNonceError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, false, nil
}
func (nm *MockNonceManagerGetChainNonceError) UpdateCacheNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerGetChainNonceError) GetLock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerGetChainNonceError) ReleaseLock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

func TestNonceHandlerErrors(t *testing.T) {
	nonceManagers := []infra.NonceManager{
		&MockNonceManagerGetNonceError{},
		&MockNonceManagerUpdateCacheNonceError{},
		&MockNonceManagerGetLockError{},
		&MockNonceManagerGetChainNonceError{},
	}

	for _, nm := range nonceManagers {
		var nonceH types.HandlerFunc
		if reflect.TypeOf(nm) == reflect.TypeOf(&MockNonceManagerGetChainNonceError{}) {
			nonceH = NonceHandler(nm, getChainNonce(true))
		} else {
			nonceH = NonceHandler(nm, getChainNonce(false))
		}

		ctx := makeNonceContext()
		nonceH(ctx)
		if len(ctx.T.Errors) == 0 {
			t.Errorf("NonceHandler: An error should have been added to context for %T", nm)
		}
	}
}

type ethClientMock struct{}

const nonce uint64 = 12

func (ec *ethClientMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return nonce, nil
}
func TestGetChainNonce(t *testing.T) {
	client := &ethClientMock{}

	cid := big.NewInt(36)
	a := common.HexToAddress("0xabcdabcdabcdabcdabcdabcd")
	getNonceFunc := GetChainNonce(client)
	n, err := getNonceFunc(cid, &a)
	if err != nil {
		t.Error("NonceHandler: error should have been nil")
	}
	if n != nonce {
		t.Errorf("NonceHandler: should have returned %v, got %v instead", nonce, n)
	}
}
