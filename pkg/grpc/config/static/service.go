package static

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type Services struct {
	Contracts *Contracts
	Envelopes *Envelopes
	Mock      *Mock
}

func (i *Services) Field() (interface{}, error) {
	return utils.ExtractField(i)
}

type Contracts struct{}

type Envelopes struct{}
