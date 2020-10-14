package ethereum

type ETHAccountResponse struct {
	Address             string `json:"address" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	PublicKey           string `json:"publicKey"`
	CompressedPublicKey string `json:"compressedPublicKey"`
	KeyType             string `json:"keyType" example:"Secp256k1"`
	Namespace           string `json:"namespace,omitempty" example:"tenant_id"`
}
