package keymanager

type SignPayloadRequest struct {
	Data      string `json:"data" validate:"required,isHex" example:"0x6d79206461746120746f207369676e"`
	Namespace string `json:"namespace,omitempty" example:"tenant_id"`
}
