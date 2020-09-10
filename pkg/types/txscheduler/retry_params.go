package txscheduler

import (
	"math"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type GasPriceParams struct {
	Priority    string      `json:"priority,omitempty" validate:"isPriority" example:"very-high"`
	RetryPolicy RetryParams `json:"retryPolicy,omitempty"`
}
type RetryParams struct {
	Interval  string  `json:"interval,omitempty" validate:"omitempty,minDuration=1s" example:"2m"`
	Increment float64 `json:"increment,omitempty" validate:"omitempty" example:"0.05"`
	Limit     float64 `json:"limit,omitempty" validate:"omitempty" example:"0.5"`
}
type IntervalRetryParams struct {
	Interval string `json:"interval,omitempty" validate:"omitempty,isDuration" example:"2m"`
}

const SentryMaxRetries = 10

func (g *RetryParams) Validate() error {
	if err := utils.GetValidator().Struct(g); err != nil {
		return err
	}

	// required_with does not work with floats as the 0 value is valid
	if g.Limit > 0 && g.Increment == 0 {
		return errors.InvalidParameterError("fields 'increment' must be specified when 'limit' is set")
	}
	if g.Increment > 0 && g.Limit == 0 {
		return errors.InvalidParameterError("field 'limit' must be specified when 'increment' is set")
	} else if g.Increment > 0 && math.Ceil(g.Limit/g.Increment) > SentryMaxRetries {
		return errors.InvalidParameterError("Maximum amount of retries %d was exceeded. Reduce 'limit' or increase 'increment` values", SentryMaxRetries)
	}

	return nil
}
