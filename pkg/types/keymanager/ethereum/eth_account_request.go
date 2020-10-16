package ethereum

type CreateETHAccountRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}

type ImportETHAccountRequest struct {
	PrivateKey string `json:"privateKey" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249" validate:"required"`
	Namespace  string `json:"namespace,omitempty" example:"tenant_id"`
}
