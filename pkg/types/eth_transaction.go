package types

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type ETHTransaction struct {
	Hash           string    `json:"hash,omitempty" example:"0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e"`
	From           string    `json:"from,omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	To             string    `json:"to,omitempty" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534"`
	Nonce          string    `json:"nonce,omitempty" example:"1"`
	Value          string    `json:"value,omitempty" example:"71500000 (wei)"`
	GasPrice       string    `json:"gasPrice,omitempty" example:"71500000 (wei)"`
	Gas            string    `json:"gas,omitempty" example:"21000"`
	Data           string    `json:"data,omitempty" example:"0xfe378324abcde723"`
	Raw            string    `json:"raw,omitempty" example:"0xfe378324abcde723"`
	PrivateFrom    string    `json:"privateFrom,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivateFor     []string  `json:"privateFor,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	PrivacyGroupID string    `json:"privacyGroupID,omitempty" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
	CreatedAt      time.Time `json:"createdAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}

func (t *ETHTransaction) GetHash() ethcommon.Hash {
	return ethcommon.HexToHash(t.Hash)
}
