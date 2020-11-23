package entities

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

type ETHTransactionParams struct {
	From            string        `json:"from,omitempty"  example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To              string        `json:"to,omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	Value           string        `json:"value,omitempty"  example:"71500000 (wei)"`
	GasPrice        string        `json:"gasPrice,omitempty" example:"71500000 (wei)"`
	Gas             string        `json:"gas,omitempty" example:"21000"`
	MethodSignature string        `json:"methodSignature,omitempty" example:"transfer(address,uint256)"`
	Args            []interface{} `json:"args,omitempty"`
	Raw             string        `json:"raw,omitempty" example:"0xfe378324abcde723"`
	ContractName    string        `json:"contractName,omitempty" example:"MyContract"`
	ContractTag     string        `json:"contractTag,omitempty" example:"v1.1.0"`
	Nonce           string        `json:"nonce,omitempty" example:"1"`
	Protocol        string        `json:"protocol,omitempty" example:"Tessera"`
	PrivateFrom     string        `json:"privateFrom,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor      []string      `json:"privateFor,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID  string        `json:"privacyGroupId,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

type PrivateETHTransactionParams struct {
	PrivateFrom    string
	PrivateFor     []string
	PrivacyGroupID string
	PrivateTxType  utils.PrivateTxType
}
