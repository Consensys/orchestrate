package keymanager

type SignPayloadRequest struct {
	Data      string `json:"data" validate:"required,isHex" example:"my data to sign"`
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}
