package handlers

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
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

func (nm *MockNonceManager) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	nm.nonce = newNonce
	return nil
}

func (nm *MockNonceManager) Lock(chainID *big.Int, a *common.Address) (string, error) {
	nm.mux.Lock()
	return "random", nil
}

func (nm *MockNonceManager) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
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

func (nm *MockNonceManagerFromChain) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	nm.nonce = newNonce
	return nil
}

func (nm *MockNonceManagerFromChain) Lock(chainID *big.Int, a *common.Address) (string, error) {
	nm.mux.Lock()
	return "random", nil
}

func (nm *MockNonceManagerFromChain) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
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
func (nm *MockNonceManagerGetNonceError) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerGetNonceError) Lock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerGetNonceError) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerSetNonceError struct{}

func (nm *MockNonceManagerSetNonceError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, true, nil
}
func (nm *MockNonceManagerSetNonceError) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return errors.New("error")
}
func (nm *MockNonceManagerSetNonceError) Lock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerSetNonceError) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerLockError struct{}

func (nm *MockNonceManagerLockError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, true, nil
}
func (nm *MockNonceManagerLockError) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerLockError) Lock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", errors.New("error")
}
func (nm *MockNonceManagerLockError) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

type MockNonceManagerGetChainNonceError struct{}

func (nm *MockNonceManagerGetChainNonceError) GetNonce(chainID *big.Int, a *common.Address) (uint64, bool, error) {
	return 0, false, nil
}
func (nm *MockNonceManagerGetChainNonceError) SetNonce(chainID *big.Int, a *common.Address, newNonce uint64) error {
	return nil
}
func (nm *MockNonceManagerGetChainNonceError) Lock(chainID *big.Int, a *common.Address) (string, error) {
	return "random", nil
}
func (nm *MockNonceManagerGetChainNonceError) Unlock(chainID *big.Int, a *common.Address, lockSig string) error {
	return nil
}

func TestNonceHandlerErrors(t *testing.T) {
	nonceManagers := []services.NonceManager{
		&MockNonceManagerGetNonceError{},
		&MockNonceManagerSetNonceError{},
		&MockNonceManagerLockError{},
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
