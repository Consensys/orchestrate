package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type MockTxSigner struct {
	t *testing.T
}

func (s *MockTxSigner) Sign(chain *types.Chain, a common.Address, tx *ethtypes.Transaction) (raw []byte, hash *common.Hash, err error) {
	if chain.ID.Text(10) == "0" {
		return []byte(``), nil, fmt.Errorf("Could not sign")
	}
	h := common.HexToHash("0xabcdef")
	return hexutil.MustDecode("0xabcdef"), &h, nil
}

func makeSignerContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	switch i % 4 {
	case 0:
		ctx.T.Chain().ID = big.NewInt(10)
		ctx.T.Tx().SetRaw(hexutil.MustDecode("0xabde4f3a"))
		h := common.HexToHash("0x12345678")
		ctx.T.Tx().SetHash(&h)
		ctx.Keys["errors"] = 0
		ctx.Keys["raw"] = "0xabde4f3a"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000012345678"
	case 1:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.T.Tx().SetRaw(hexutil.MustDecode("0xabde4f3a"))
		h := common.HexToHash("0x12345678")
		ctx.T.Tx().SetHash(&h)
		ctx.Keys["errors"] = 0
		ctx.Keys["raw"] = "0xabde4f3a"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000012345678"
	case 2:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.T.Tx().SetRaw([]byte(``))
		ctx.Keys["errors"] = 1
		ctx.Keys["raw"] = "0x"
		ctx.Keys["hash"] = "0x0000000000000000000000000000000000000000000000000000000000000000"
	case 3:
		ctx.T.Chain().ID = big.NewInt(10)
		ctx.T.Tx().SetRaw([]byte(``))
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
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeSignerContext(i)
		go func(ctx *types.Context) {
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

		if hexutil.Encode(out.T.Tx().Raw()) != raw {
			t.Errorf("Signer: expected Raw %v but got %v", raw, hexutil.Encode(out.T.Tx().Raw()))
		}

		if out.T.Tx().Hash().Hex() != hash {
			t.Errorf("Signer: expected hash %v but got %v", hash, out.T.Tx().Hash().Hex())
		}
	}
}
