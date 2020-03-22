// +build unit

package signature

import (
	"testing"

	csutils "github.com/ConsenSys/golang-utils/ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	key, _ := crypto.GenerateKey()

	hash := ethcommon.HexToHash("")
	sig, err := EthECDSA.Sign(hash.Bytes(), key)
	assert.NoError(t, err, "Sign should not error")

	addr, err := csutils.EcRecover(hash, sig)
	assert.NoError(t, err, "EcRecover should not error")
	expected := crypto.PubkeyToAddress(key.PublicKey)
	assert.Equal(t, expected.Hex(), addr.Hex(), "ECRecover should return correct address")
}
