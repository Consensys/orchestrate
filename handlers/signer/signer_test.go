package signer

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type MockTxSigner struct {
	t *testing.T
}

func (s *MockTxSigner) SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
	if chain.ID().String() == "0" {
		return []byte(``), nil, fmt.Errorf("could not sign")
	}
	h := ethcommon.HexToHash("0xabcdef")
	return hexutil.MustDecode("0xabcdef"), &h, nil
}

func (s *MockTxSigner) SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) GenerateWallet() (add *ethcommon.Address, err error) {
	return nil, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error) {
	return []byte{}, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) ImportPrivateKey(priv string) (err error) {
	return fmt.Errorf("importPrivateKey not implemented")
}

func makeSignerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 4 {
	case 0:
		h := ethcommon.HexToHash("0x12345678")
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		txctx.Envelope.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a").SetHash(h)
		txctx.Set("errors", 0)
		txctx.Set("raw", "0xabde4f3a")
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000012345678")
	case 1:
		h := ethcommon.HexToHash("0x12345678")
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		txctx.Envelope.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a").SetHash(h)

		txctx.Set("errors", 0)
		txctx.Set("raw", "0xabde4f3a")
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000012345678")
	case 2:
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		txctx.Envelope.Tx = (&ethereum.Transaction{}).SetRaw(``)
		txctx.Set("errors", 1)
		txctx.Set("raw", "")
		txctx.Set("hash", "")
	case 3:
		txctx.Envelope.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		txctx.Envelope.Tx = (&ethereum.Transaction{}).SetRaw(``)
		txctx.Set("errors", 0)
		txctx.Set("raw", "0xabcdef")
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000000abcdef")
	}
	return txctx
}

func TestSigner(t *testing.T) {
	s := MockTxSigner{t: t}
	signer := Signer(&s)

	rounds := 100
	outs := make(chan *engine.TxContext, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		txctx := makeSignerContext(i)
		go func(txctx *engine.TxContext) {
			defer wg.Done()
			signer(txctx)
			outs <- txctx
		}(txctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Signer: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount, raw, hash := out.Get("errors").(int), out.Get("raw").(string), out.Get("hash").(string)
		if len(out.Envelope.Errors) != errCount {
			t.Errorf("Signer: expected %v errors but got %v", errCount, out.Envelope.Errors)
		}

		if out.Envelope.Tx.GetRaw() != raw {
			t.Errorf("Signer: expected Raw %v but got %v", raw, out.Envelope.Tx.GetRaw())
		}

		if out.Envelope.Tx.GetHash() != hash {
			t.Errorf("Signer: expected hash %v but got %v", hash, out.Envelope.Tx.GetHash())
		}
	}
}
