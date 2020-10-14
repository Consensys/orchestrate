package ethereum

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type CreateETHAccountRequest struct {
	KeyType   utils.ETHKeyType `json:"keyType" example:"Secp256k1" validate:"required,isKeyType"`
	Namespace string           `json:"namespace,omitempty" example:"tenant_id"`
}
