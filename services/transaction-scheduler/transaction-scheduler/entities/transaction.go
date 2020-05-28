package entities

import "time"

type ETHTransaction struct {
	Hash           string    `json:"hash,omitempty"`
	From           string    `json:"from,omitempty"`
	To             string    `json:"to,omitempty"`
	Nonce          string    `json:"nonce,omitempty"`
	Value          string    `json:"value,omitempty"`
	GasPrice       string    `json:"gasPrice,omitempty"`
	GasLimit       string    `json:"gasLimit,omitempty"`
	Data           string    `json:"data,omitempty"`
	Raw            string    `json:"raw,omitempty"`
	PrivateFrom    string    `json:"privateFrom,omitempty"`
	PrivateFor     []string  `json:"privateFor,omitempty"`
	PrivacyGroupID string    `json:"privacyGroupID,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}
