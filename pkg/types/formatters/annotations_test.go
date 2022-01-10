// +build unit

package formatters

import (
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/assert"
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
