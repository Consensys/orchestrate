package api

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"math/big"
)

// SerializedAddress is a helper for address serialization
type SerializedAddress string

// ToGeth deserialize an hex encoded ethereum address
func (a *SerializedAddress) ToGeth() (common.Address) {
	return common.HexToAddress(string(*a))
}


// SerializedTx maps the request body with an Ethereum Transaction
type SerializedTx struct {
	Nonce uint64 `json:"nonce" binding:"required"`
	To string `json:"to" binding:"required"`
	Value string `json:"value" binding:"required"`
	GasLimit uint64 `json:"gaslimit" binding:"required"`
	GasPrice string `json:"gasprice" binding:"required"`
	Data string `json:"data" binding:"required"`
}

// ToGeth converts the serialized tx to a geth tx
func (t *SerializedTx) ToGeth() (g *ethtypes.Transaction) {
	return ethtypes.NewTransaction(
		t.getNonce(),
		t.getTo(),
		t.getValue(),
		t.getGasLimit(),
		t.getGasPrice(),
		t.getData(),
	)
}

func (t *SerializedTx) getNonce() uint64 {
	return t.Nonce
}

func (t *SerializedTx) getTo() (a common.Address) {
	return common.HexToAddress(t.To)
}

func (t *SerializedTx) getValue() (*big.Int) {
	i := new(big.Int)
	i.SetString(t.Value, 10)
	return i
}

func (t *SerializedTx) getGasLimit() uint64 {
	return t.GasLimit
}

func (t *SerializedTx) getGasPrice() (*big.Int) {
	i := new(big.Int)
	i.SetString(t.GasPrice, 10)
	return i
}

func (t *SerializedTx) getData() []byte {
	return []byte(t.Data)
}

// JsonifiedChain constaines a Jsonified Chain object
type JsonifiedChain struct {
	ID string `json:"id" binding:"required"`
	IsEIP155 bool `json:"isEIP155" binding:"required"`
}

// ToCoreStack converts the JSON object to its actual form
func (c *JsonifiedChain) ToCoreStack() (*types.Chain) {
	bigID := new(big.Int)
	bigID.SetString(c.ID, 16)
	return &types.Chain{
		ID: bigID,
		IsEIP155: c.IsEIP155,
	}
}


