package ethereum

type CreateETHAccountRequest struct {
	KeyType   string `json:"keyType" example:"Secp256k1" validate:"required,isKeyType"`
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}
