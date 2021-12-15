package api

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CreateAccountRequest struct {
	Alias      string            `json:"alias" validate:"omitempty" example:"personal-account" `
	Chain      string            `json:"chain" validate:"omitempty" example:"besu"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type ImportAccountRequest struct {
	Alias      string            `json:"alias" validate:"omitempty" example:"personal-account"`
	Chain      string            `json:"chain" validate:"omitempty" example:"quorum"`
	PrivateKey hexutil.Bytes     `json:"privateKey" validate:"required" example:"0x66232652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2D" swaggertype:"string"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type UpdateAccountRequest struct {
	Alias      string            `json:"alias" validate:"omitempty"  example:"personal-account"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
