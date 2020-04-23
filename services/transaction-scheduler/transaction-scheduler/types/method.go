package types

const (
	MethodSendRawTransaction        = "ETH_SENDRAWTRANSACTION"        // Classic ETH
	MethodSendPrivateTransaction    = "ETH_SENDPRIVATETRANSACTION"    // Quorum Constellation
	MethodSendRawPrivateTransaction = "ETH_SENDRAWPRIVATETRANSACTION" // Quorum Tessera
	MethodEEASendPrivateTransaction = "EEA_SENDPRIVATETRANSACTION"    // Besu Orion
)
