package txscheduler

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FormatAnnotationsToInternalData(annotations Annotations, parentJobUUID string) *entities.InternalData {
	internalData := &entities.InternalData{
		OneTimeKey:        annotations.OneTimeKey,
		Priority:          annotations.GasPricePolicy.Priority,
		GasPriceIncrement: annotations.GasPricePolicy.RetryPolicy.Increment,
		GasPriceLimit:     annotations.GasPricePolicy.RetryPolicy.Limit,
		HasBeenRetried:    annotations.HasBeenRetried,
	}

	if parentJobUUID != "" {
		internalData.ParentJobUUID = parentJobUUID
	}

	if annotations.GasPricePolicy.RetryPolicy.Interval != "" {
		// we can skip the error check as at this point we know the interval is a duration as it already passed validation
		internalData.RetryInterval, _ = time.ParseDuration(annotations.GasPricePolicy.RetryPolicy.Interval)
	}

	if internalData.Priority == "" {
		internalData.Priority = utils.PriorityMedium
	}

	return internalData
}

func FormatInternalDataToAnnotations(data *entities.InternalData) Annotations {
	gasPricePolicy := GasPriceParams{
		Priority: data.Priority,
		RetryPolicy: RetryParams{
			Increment: data.GasPriceIncrement,
			Limit:     data.GasPriceLimit,
		},
	}

	if data.RetryInterval != 0 {
		gasPricePolicy.RetryPolicy.Interval = data.RetryInterval.String()
	}

	return Annotations{
		OneTimeKey:     data.OneTimeKey,
		GasPricePolicy: gasPricePolicy,
		HasBeenRetried: data.HasBeenRetried,
	}
}
