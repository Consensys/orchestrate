package zksnarks

type VerifyPayloadRequest struct {
	Data      string `json:"data" validate:"required" example:"my data to sign"`
	Signature string `json:"signature" validate:"required" example:"0x6019a3c8..."`
	PublicKey string `json:"publicKey" validate:"required" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
}
