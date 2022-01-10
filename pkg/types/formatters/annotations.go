package formatters

import (
	"time"

	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

func FormatAnnotationsToInternalData(annotations api.Annotations) *entities.InternalData {
	internalData := &entities.InternalData{
		OneTimeKey:        annotations.OneTimeKey,
		Priority:          annotations.GasPricePolicy.Priority,
		GasPriceIncrement: annotations.GasPricePolicy.RetryPolicy.Increment,
		GasPriceLimit:     annotations.GasPricePolicy.RetryPolicy.Limit,
		HasBeenRetried:    annotations.HasBeenRetried,
	}

	if annotations.GasPricePolicy.RetryPolicy.Interval != "" {
		// we can skip the error check as at this point we know the interval is a duration as it already passed validation
		internalData.RetryInterval, _ = time.ParseDuration(annotations.GasPricePolicy.RetryPolicy.Interval)
	}

	return internalData
}

func FormatInternalDataToAnnotations(data *entities.InternalData) api.Annotations {
	gasPricePolicy := api.GasPriceParams{
		Priority: data.Priority,
		RetryPolicy: api.RetryParams{
			Increment: data.GasPriceIncrement,
			Limit:     data.GasPriceLimit,
		},
	}

	if data.RetryInterval != 0 {
		gasPricePolicy.RetryPolicy.Interval = data.RetryInterval.String()
	}

	return api.Annotations{
		OneTimeKey:     data.OneTimeKey,
		GasPricePolicy: gasPricePolicy,
		HasBeenRetried: data.HasBeenRetried,
	}
}
