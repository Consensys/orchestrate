package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
)

type MockTxSigner struct {
	t *testing.T
}

<<<<<<< HEAD
func (s *MockTxSigner) SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
=======
func (s *MockTxSigner) Sign(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
>>>>>>> master
	if chain.ID().String() == "0" {
		return []byte(``), nil, fmt.Errorf("Could not sign")
	}
	h := ethcommon.HexToHash("0xabcdef")
	return hexutil.MustDecode("0xabcdef"), &h, nil
}

func (s *MockTxSigner) SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, fmt.Errorf("SignMsg not implemented")
}

func (s *MockTxSigner) GenerateWallet() (add *ethcommon.Address, err error) {
	return nil, fmt.Errorf("SignMsg not implemented")
}

func (s *MockTxSigner) SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error) {
	return []byte{}, fmt.Errorf("SignMsg not implemented")
}

func makeSignerContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 4 {
	case 0:
		h := ethcommon.HexToHash("0x12345678")
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a").SetHash(h)
		ctx.Keys["errors"] = 0
		ctx.Keys["raw"] = "0xabde4f3a"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000012345678"
	case 1:
		h := ethcommon.HexToHash("0x12345678")
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw("0xabde4f3a").SetHash(h)

		ctx.Keys["errors"] = 0
		ctx.Keys["raw"] = "0xabde4f3a"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000012345678"
	case 2:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.Keys["errors"] = 1
		ctx.Keys["raw"] = ""
		ctx.Keys["hash"] = ""
	case 3:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(10))
		ctx.T.Tx = (&ethereum.Transaction{}).SetRaw(``)
		ctx.Keys["errors"] = 0
		ctx.Keys["raw"] = "0xabcdef"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000000abcdef"
	}
	return ctx
}

func TestSigner(t *testing.T) {
	s := MockTxSigner{t: t}
	signer := Signer(&s)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeSignerContext(i)
		go func(ctx *worker.Context) {
			defer wg.Done()
			signer(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Signer: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount, raw, hash := out.Keys["errors"].(int), out.Keys["raw"].(string), out.Keys["hash"].(string)
		if len(out.T.Errors) != errCount {
			t.Errorf("Signer: expected %v errors but got %v", errCount, out.T.Errors)
		}

		if out.T.Tx.GetRaw() != raw {
			t.Errorf("Signer: expected Raw %v but got %v", raw, out.T.Tx.GetRaw())
		}

		if out.T.Tx.GetHash() != hash {
			t.Errorf("Signer: expected hash %v but got %v", hash, out.T.Tx.GetHash())
		}
	}
}
