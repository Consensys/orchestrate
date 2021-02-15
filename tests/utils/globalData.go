package utils

type TestData struct {
	Nodes TestDataNodes `json:"nodes"`
}

type TestDataNodes struct {
	Besu   []TestDataChain `json:"besu,omitempty"`
	Quorum []TestDataChain `json:"quorum,omitempty"`
	Geth   []TestDataChain `json:"geth,omitempty"`
}

type TestDataChain struct {
	URLs              []string                 `json:"URLs,omitempty"`
	PrivateAddress    []string                 `json:"privateAddress,omitempty"`
	FundedPublicKeys  []string                 `json:"fundedPublicKeys,omitempty"`
	FundedPrivateKeys []string                 `json:"fundedPrivateKeys,omitempty"`
	PrivateTxManager  TestDataPrivateTxManager `json:"privateTxManager,omitempty"`
}

type TestDataPrivateTxManager struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}
