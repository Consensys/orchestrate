// +build unit

package dynamic_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic/testutils"
)

func TestMerge(t *testing.T) {
	mergedCfg := dynamic.Merge(testutils.Configs).(*dynamic.Configuration)
	assert.True(t, reflect.DeepEqual(mergedCfg, testutils.Config))
}
