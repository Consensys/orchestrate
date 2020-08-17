package txschedulertypes

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type Annotations struct {
	OneTimeKey  bool                `json:"oneTimeKey,omitempty" example:"true"`
	Priority    string              `json:"priority,omitempty" validate:"isPriority" example:"very-high"`
	RetryPolicy GasPriceRetryParams `json:"gasPriceRetryPolicy,omitempty"`
}

func (g *Annotations) Validate() error {
	if err := utils.GetValidator().Struct(g); err != nil {
		return err
	}

	if err := g.RetryPolicy.Validate(); err != nil {
		return err
	}

	return nil
}
