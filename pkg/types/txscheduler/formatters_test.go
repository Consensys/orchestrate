package txscheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func TestInternalLabelsAnnotationFormatters(t *testing.T) {
	internal := &entities.InternalData{
		Priority:          utils.PriorityMedium,
		RetryInterval:     5 * time.Second,
		HasBeenRetried:    true,
		OneTimeKey:        true,
		GasPriceIncrement: 0.3,
		GasPriceLimit:     0.6,
	}

	annotations := FormatInternalDataToAnnotations(internal)
	finalInternal := FormatAnnotationsToInternalData(annotations)
	assert.Equal(t, finalInternal, internal)
}
