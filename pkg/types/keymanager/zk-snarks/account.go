package zksnarks

type CreateZKSAccountRequest struct {
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}

type ZKSAccountResponse struct {
	Curve            string `json:"curve" example:"bn256"`
	SigningAlgorithm string `json:"signingAlgorithm" example:"eddsa"`
	PublicKey        string `json:"publicKey" example:"20199690451585786844338768304582194735444460424798515739606133903768949456887"`
	Namespace        string `json:"namespace,omitempty" example:"tenant_id"`
}
