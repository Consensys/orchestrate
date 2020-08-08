package types

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type IntervalRetryParams struct {
	Interval string `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
}
type GasPriceRetryParams struct {
	Interval               string  `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
	GasPriceIncrementLevel string  `json:"gasPriceIncrementLevel,omitempty" validate:"omitempty,isGasIncrementLevel" example:"medium"`
	GasPriceIncrement      float64 `json:"gasPriceIncrement,omitempty" validate:"omitempty,eqfield=GasPriceLimit|ltcsfield=GasPriceLimit" example:"1.05"`
	GasPriceLimit          float64 `json:"gasPriceLimit,omitempty" validate:"required_with=GasPriceIncrementLevel GasPriceIncrement,omitempty" example:"1.2"`
}

func (g *GasPriceRetryParams) Validate() error {
	if err := utils.GetValidator().Struct(g); err != nil {
		return err
	}

	if g.GasPriceIncrement > 0 && g.GasPriceIncrementLevel != "" {
		return errors.InvalidParameterError("fields 'gasPriceIncrement' and 'gasPriceIncrementLevel' are mutually exclusive")
	}

	if (g.Interval != "" || g.GasPriceLimit > 0) && (g.GasPriceIncrement == 0 && g.GasPriceIncrementLevel == "") {
		return errors.InvalidParameterError("fields 'gasPriceIncrement' and 'gasPriceIncrementLevel' cannot be both empty")
	}

	if (g.Interval != "" && g.GasPriceLimit == 0) || (g.Interval == "" && g.GasPriceLimit > 0) {
		return errors.InvalidParameterError("fields 'Interval' and 'GasPriceLimit' are both required")
	}

	return nil
}
