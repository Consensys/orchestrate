package utils

import "github.com/ethereum/go-ethereum/common/hexutil"

type TestData struct {
	Nodes TestDataNodes `json:"nodes"`
	OIDC  *TestOIDC     `json:"oidc"`
}

type TestDataNodes struct {
	Besu   []TestDataChain `json:"besu,omitempty"`
	Quorum []TestDataChain `json:"quorum,omitempty"`
	Geth   []TestDataChain `json:"geth,omitempty"`
}

type TestOIDC struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	TokenURL     string `json:"tokenURL"`
}

type TestDataChain struct {
	URLs              []string                 `json:"URLs,omitempty"`
	PrivateAddress    []string                 `json:"privateAddress,omitempty"`
	FundedPublicKeys  []string                 `json:"fundedPublicKeys,omitempty"`
	FundedPrivateKeys []hexutil.Bytes          `json:"fundedPrivateKeys,omitempty"`
	PrivateTxManager  TestDataPrivateTxManager `json:"privateTxManager,omitempty"`
}

type TestDataPrivateTxManager struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}
