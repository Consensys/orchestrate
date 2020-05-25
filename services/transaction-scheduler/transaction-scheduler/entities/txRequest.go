package entities

import (
	"time"
)

const (
	MethodSendRawTransaction        = "ETH_SENDRAWTRANSACTION"        // Classic ETH
	MethodSendPrivateTransaction    = "ETH_SENDPRIVATETRANSACTION"    // Quorum Constellation
	MethodSendRawPrivateTransaction = "ETH_SENDRAWPRIVATETRANSACTION" // Quorum Tessera
	MethodEEASendPrivateTransaction = "EEA_SENDPRIVATETRANSACTION"    // Besu Orion
)

type TxRequest struct {
	IdempotencyKey string
	Schedule       *Schedule
	Params         *TxRequestParams
	Labels         map[string]string
	CreatedAt      time.Time
}

type TxRequestParams struct {
	From            *string
	To              *string
	Value           *string
	GasPrice        *string
	GasLimit        *string
	MethodSignature *string
	Args            []string
	Raw             *string
	ContractName    *string
	ContractTag     *string
}
