package signer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/magiconair/properties/assert"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	eeaHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/eea"
	ethereumHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/ethereum"
	tesseraHandlers "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

type MockTxSigner struct {
	t *testing.T
}

type MockTesseraClient struct {
	t *testing.T
}

var alreadySignedTx = "0x04"
var signedTx = "0x01"
var signedPrivateTx = "0x02"
var signedTesseraTx = "0x03"

func (s *MockTxSigner) SignTx(_ context.Context, netChain *big.Int, _ ethcommon.Address, _ *ethtypes.Transaction) (raw []byte, hash *ethcommon.Hash, err error) {
	if netChain.String() == "0" {
		return []byte(``), nil, fmt.Errorf("could not sign public ethereum transaction")
	}
	h := ethcommon.HexToHash("0xabcdef")
	return hexutil.MustDecode(signedTx), &h, nil
}

func (s *MockTxSigner) SignPrivateEEATx(_ context.Context, netChain *big.Int, _ ethcommon.Address, _ *ethtypes.Transaction, _ *types.PrivateArgs) (raw []byte, hash *ethcommon.Hash, err error) {
	if netChain.String() == "0" {
		return []byte(``), nil, fmt.Errorf("could not sign eea transaction")
	}
	h := ethcommon.HexToHash("0xabcdef")
	return hexutil.MustDecode(signedPrivateTx), &h, nil
}

func (s *MockTxSigner) SignPrivateTesseraTx(_ context.Context, netChain *big.Int, _ ethcommon.Address, _ *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error) {
	if netChain.String() == "0" {
		return []byte(``), nil, fmt.Errorf("could not sign tessera transaction")
	}
	h := ethcommon.HexToHash("0xabcdef")
	return hexutil.MustDecode(signedTesseraTx), &h, nil
}

func (s *MockTxSigner) SignMsg(_ context.Context, _ ethcommon.Address, _ string) (rsv []byte, hash *ethcommon.Hash, err error) {
	return []byte{}, nil, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) GenerateAccount(_ context.Context) (add *ethcommon.Address, err error) {
	return nil, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) SignRawHash(_ ethcommon.Address, _ []byte) (rsv []byte, err error) {
	return []byte{}, fmt.Errorf("signMsg not implemented")
}

func (s *MockTxSigner) ImportPrivateKey(_ context.Context, _ string) (err error) {
	return fmt.Errorf("importPrivateKey not implemented")
}

func (tc *MockTesseraClient) AddClient(_ string, _ tessera.EnclaveEndpoint) {}

func (tc *MockTesseraClient) StoreRaw(chainID string, _ []byte, _ string) (txHash []byte, err error) {
	if chainID == "0" {
		return []byte(``), fmt.Errorf("mock: store raw failed")
	}
	return hexutil.MustDecode("0xabcdef"), nil
}

func (tc *MockTesseraClient) GetStatus(chainID string) (status string, err error) {
	if chainID == "0" {
		return "", fmt.Errorf("mock: get status failed")
	}
	return "", nil
}

func makeSignerContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())

	switch i % 1 {
	case 0:
		_ = txctx.Envelope.
			SetChainIDUint64(10).
			SetGas(10).
			SetNonce(11).
			SetGasPrice(big.NewInt(12)).
			MustSetToString("0x1").
			MustSetFromString("0x2").
			MustSetTxHashString("0x12345678").
			SetRawString(alreadySignedTx)
		txctx.Set("errors", 0)
		txctx.Set("raw", alreadySignedTx)
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000012345678")
	case 1:
		_ = txctx.Envelope.
			SetChainIDUint64(0).
			SetGas(10).
			SetNonce(11).
			SetGasPrice(big.NewInt(12)).
			MustSetToString("0x0").
			MustSetFromString("0x0").
			SetTxHash(ethcommon.HexToHash("0x12345678")).
			SetRawString(alreadySignedTx)
		txctx.Set("errors", 0)
		txctx.Set("raw", alreadySignedTx)
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000012345678")
	case 2:
		_ = txctx.Envelope.SetChainIDUint64(0)
		txctx.Set("errors", 1)
		txctx.Set("raw", "")
		txctx.Set("hash", "")
	case 3:
		_ = txctx.Envelope.SetChainIDUint64(10)
		txctx.Set("errors", 0)
		txctx.Set("raw", signedTx)
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000000abcdef")
	case 4:
		_ = txctx.Envelope.SetChainIDUint64(10).SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).SetDataString("")
		txctx.Set("errors", 0)
		txctx.Set("errors", 0)
		txctx.Set("raw", signedTesseraTx)
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000000abcdef")
	case 5:
		_ = txctx.Envelope.SetChainIDUint64(10).SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION)
		txctx.Set("errors", 1)
		txctx.Set("raw", "")
		txctx.Set("hash", "")
	case 6:
		_ = txctx.Envelope.SetChainIDUint64(0).SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).SetDataString("")
		txctx.Set("errors", 1)
		txctx.Set("raw", "")
		txctx.Set("hash", "")
	case 7:
		_ = txctx.Envelope.SetChainIDUint64(10).SetMethod(tx.Method_EEA_SENDPRIVATETRANSACTION).SetDataString("")
		txctx.Set("errors", 0)
		txctx.Set("raw", signedPrivateTx)
		txctx.Set("hash", "0x0000000000000000000000000000000000000000000000000000000000abcdef")
	case 8:
		txctx.Set("errors", 1)
		txctx.Set("raw", "")
		txctx.Set("hash", "")
	}
	return txctx
}

func TestSigner(t *testing.T) {

	s := &MockTxSigner{t: t}
	tc := &MockTesseraClient{t: t}

	signer := TxSigner(
		eeaHandlers.Signer(s),
		ethereumHandlers.Signer(s),
		tesseraHandlers.Signer(s, tc),
	)

	rounds := 25
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
		assert.Equal(t, len(out.Envelope.Errors), errCount, fmt.Sprintf("Signer: expected %v errors but got %v", errCount, out.Envelope.Errors))
		assert.Equal(t, out.Envelope.GetRaw(), raw, fmt.Sprintf("Signer: expected Raw %v but got %v", raw, out.Envelope.GetRaw()))
		assert.Equal(t, out.Envelope.GetTxHash().Hex(), hash, fmt.Sprintf("Signer: expected hash %v but got %v", hash, out.Envelope.MustGetTxHashValue().Hex()))
	}
}
