package entities

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ETHAccount struct {
	Address             ethcommon.Address `json:"address" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	PublicKey           hexutil.Bytes     `json:"publicKey"`
	CompressedPublicKey hexutil.Bytes     `json:"compressedPublicKey"`
	Namespace           string            `json:"namespace,omitempty" example:"tenant_id"`
}
