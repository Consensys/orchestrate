package entities

type FaucetRequest struct {
	Chain       *Chain
	Beneficiary string
	Candidates  map[string]*Faucet
}
