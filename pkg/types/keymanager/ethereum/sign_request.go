package ethereum

type SignETHTransactionRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
	Nonce     uint64 `json:"nonce" example:"1"`
	To        string `json:"to" validate:"required,isHex" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	Amount    string `json:"amount,omitempty" validate:"isBig" example:"100000000000"`
	GasPrice  string `json:"gasPrice" validate:"required,isBig" example:"100000000000"`
	GasLimit  uint64 `json:"gasLimit" validate:"required" example:"21000"`
	Data      string `json:"data,omitempty" validate:"isHex" example:"0xfeaeee..."`
	ChainID   string `json:"chainID" validate:"required,isBig" example:"1 (mainnet)"`
}
