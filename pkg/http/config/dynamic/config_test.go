// +build unit

package dynamic_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic/testutils"
)

func TestMerge(t *testing.T) {
	mergedCfg := dynamic.Merge(testutils.Configs).(*dynamic.Configuration)
	assert.True(t, reflect.DeepEqual(mergedCfg, testutils.Config))
}
