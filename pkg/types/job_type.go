package types

const (
	EthereumTransaction       = "eth://ethereum/transaction"       // Classic public Ethereum transaction
	EthereumRawTransaction    = "eth://ethereum/rawTransaction"    // Classic raw transaction
	OrionMarkingTransaction   = "eth://orion/markingTransaction"   // Besu public transaction
	OrionEEATransaction       = "eth://orion/eeaTransaction"       // Besu private tx for Orion
	TesseraPublicTransaction  = "eth://tessera/publicTransaction"  // Tessera public transaction
	TesseraPrivateTransaction = "eth://tessera/privateTransaction" // Tessera private transaction
)
