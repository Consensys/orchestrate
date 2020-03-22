// +build unit

package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	msg := &Message{Name: "provider-test", Conf: "conf-test"}
	assert.Equal(t, "provider-test", msg.ProviderName())
	assert.Equal(t, "conf-test", msg.Configuration().(string))
}

func TestMerge(t *testing.T) {
	confs := map[string]interface{}{
		"provider1": "conf1",
		"provider2": "conf2",
	}

	merged := MergeConfiguration(confs).([]string)
	assert.Equal(t, []string{"conf1@provider1", "conf2@provider2"}, merged, "Merge should be correct")
}
