package ethereum

type ETHAccountResponse struct {
	Address             string `json:"address" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	PublicKey           string `json:"publicKey" example:"0x047c25564c1b6a1553fa8f204be4229439e02b728ca28697003dc1c96ae51ff2c4d686e8494b3c1aeab21d7c3e88f0e0b418744e3bfb747581e8a68a97503add03"`
	CompressedPublicKey string `json:"compressedPublicKey" example:"0x037c25564c1b6a1553fa8f204be4229439e02b728ca28697003dc1c96ae51ff2c4"`
	Namespace           string `json:"namespace,omitempty" example:"tenant_id"`
}

type CreateETHAccountRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}

type ImportETHAccountRequest struct {
	PrivateKey string `json:"privateKey" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249" validate:"required"`
	Namespace  string `json:"namespace,omitempty" example:"tenant_id"`
}
