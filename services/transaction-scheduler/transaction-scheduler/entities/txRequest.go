package entities

import (
	"time"
)

type TxRequest struct {
	IdempotencyKey string
	Schedule       *Schedule
	Params         *TxRequestParams
	Labels         map[string]string
	CreatedAt      time.Time
}

type TxRequestParams struct {
	From            string   `json:"from,omitempty"`
	To              string   `json:"to,omitempty"`
	Value           string   `json:"value,omitempty"`
	GasPrice        string   `json:"gasPrice,omitempty"`
	GasLimit        string   `json:"gas,omitempty"`
	MethodSignature string   `json:"methodSignature,omitempty"`
	Args            []string `json:"args,omitempty"`
	Raw             string   `json:"raw,omitempty"`
	ContractName    string   `json:"contractName,omitempty"`
	ContractTag     string   `json:"contractTag,omitempty"`
	Nonce           string   `json:"nonce,omitempty"`
	PrivateFrom     string   `json:"privateFrom,omitempty"`
	PrivateFor      []string `json:"privateFor,omitempty"`
	PrivacyGroupID  string   `json:"privacyGroupdID,omitempty"`
}
