package mock

import (
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/handler"
)

type Handler struct {
	mux *sync.Mutex

	SetupCalls   int32
	CleanupCalls int32

	GetInitialPositionCalls map[string]int32
	ListenCalls             map[string]int32

	Receipts map[string]int32
	Errors   map[string]int32
	Blocks   map[string]int32
}

func NewHandler() *Handler {
	return &Handler{
		mux:                     &sync.Mutex{},
		GetInitialPositionCalls: make(map[string]int32),
		ListenCalls:             make(map[string]int32),
		Receipts:                make(map[string]int32),
		Errors:                  make(map[string]int32),
		Blocks:                  make(map[string]int32),
	}
}

func (h *Handler) Setup(session handler.TxListenerSession) error {
	log.Warn("setup")

	h.mux.Lock()
	defer h.mux.Unlock()

	h.SetupCalls++

	for _, chain := range session.Chains() {
		h.GetInitialPositionCalls[chain.Text(10)] = 0
		h.ListenCalls[chain.Text(10)] = 0
		h.Receipts[chain.Text(10)] = 0
		h.Errors[chain.Text(10)] = 0
		h.Blocks[chain.Text(10)] = 0
	}

	return nil
}

func (h *Handler) GetInitialPosition(chain *big.Int) (blockNumber, txIndex int64, err error) {
	log.Warnf("getInitialPosition %q", chain.Text(10))

	h.mux.Lock()
	defer h.mux.Unlock()

	h.GetInitialPositionCalls[chain.Text(10)]++
	return 0, 0, nil
}

func (h *Handler) Cleanup(session handler.TxListenerSession) error {
	log.Warn("cleanup")

	h.mux.Lock()
	defer h.mux.Unlock()

	h.CleanupCalls++
	return nil
}

func (h *Handler) Listen(session handler.TxListenerSession, l handler.ChainListener) error {
	h.mux.Lock()
	h.ListenCalls[l.ChainID().Text(10)]++
	h.mux.Unlock()

	wait := &sync.WaitGroup{}
	wait.Add(3)

	go func() {
	blockLoop:
		for {
			select {
			case <-l.Context().Done():
				break blockLoop
			case _, ok := <-l.Blocks():
				if !ok {
					break blockLoop
				} else {
					h.mux.Lock()
					h.Blocks[l.ChainID().Text(10)]++
					h.mux.Unlock()
				}
			}
		}
		wait.Done()
	}()

	go func() {
	receiptLoop:
		for {
			select {
			case <-l.Context().Done():
				break receiptLoop
			case _, ok := <-l.Receipts():
				if !ok {
					break receiptLoop
				} else {
					h.mux.Lock()
					h.Receipts[l.ChainID().Text(10)]++
					h.mux.Unlock()
				}
			}
		}
		wait.Done()
	}()

	go func() {
	errorLoop:
		for {
			select {
			case <-l.Context().Done():
				break errorLoop
			case _, ok := <-l.Errors():
				if !ok {
					break errorLoop
				} else {
					h.mux.Lock()
					h.Errors[l.ChainID().Text(10)]++
					h.mux.Unlock()
				}
			}
		}
		wait.Done()
	}()

	wait.Wait()
	return nil
}
