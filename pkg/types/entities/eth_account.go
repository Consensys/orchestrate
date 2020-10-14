package entities

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type ETHAccount struct {
	Address             string           `json:"address" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	PublicKey           string           `json:"publicKey"`
	CompressedPublicKey string           `json:"compressedPublicKey"`
	Namespace           string           `json:"namespace,omitempty" example:"tenant_id"`
	KeyType             utils.ETHKeyType `json:"keyType" example:"Secp256k1"`
}
