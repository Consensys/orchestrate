package ethereum

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type CreateETHAccountRequest struct {
	KeyType   utils.ETHKeyType `json:"keyType" example:"Secp256k1" validate:"required,isKeyType"`
	Namespace string           `json:"namespace,omitempty" example:"tenant_id"`
}

type ImportETHAccountRequest struct {
	PrivateKey string `json:"privateKey" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249" validate:"required"`
	KeyType    string `json:"keyType" example:"Secp256k1" validate:"required,isKeyType"`
	Namespace  string `json:"namespace,omitempty" example:"tenant_id"`
}
