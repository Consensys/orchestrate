package txschedulertypes

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type IntervalRetryParams struct {
	Interval string `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
}
type GasPriceRetryParams struct {
	Interval       string  `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
	IncrementLevel string  `json:"incrementLevel,omitempty" validate:"omitempty,isGasIncrementLevel" example:"medium"`
	Increment      float64 `json:"increment,omitempty" validate:"omitempty,eqfield=Limit|ltcsfield=Limit" example:"1.05"`
	Limit          float64 `json:"limit,omitempty" validate:"omitempty" example:"1.2"`
}

func (g *GasPriceRetryParams) Validate() error {
	if err := utils.GetValidator().Struct(g); err != nil {
		return err
	}

	if g.Increment > 0 && g.IncrementLevel != "" {
		return errors.InvalidParameterError("fields 'increment' and 'incrementLevel' are mutually exclusive")
	}

	// required_with does not work with floats as the 0 value is valid
	if g.Limit > 0 && (g.IncrementLevel == "" && g.Increment == 0) {
		return errors.InvalidParameterError("fields 'increment' or 'incrementLevel' must be specified when 'limit' is set")
	}
	if g.Increment > 0 && g.Limit == 0 {
		return errors.InvalidParameterError("field 'limit' must be specified when 'increment' is set")
	}
	if g.IncrementLevel != "" && g.Limit == 0 {
		return errors.InvalidParameterError("field 'limit' must be specified when 'incrementLevel' is set")
	}

	return nil
}
