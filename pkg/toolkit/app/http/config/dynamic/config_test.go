// +build unit

package dynamic_test

import (
	"reflect"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	mergedCfg := dynamic.Merge(testutils.Configs).(*dynamic.Configuration)
	assert.True(t, reflect.DeepEqual(mergedCfg, testutils.Config))
}
