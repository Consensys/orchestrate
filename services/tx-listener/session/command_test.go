// +build unit

package session

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
)

var config1 = &dynamic.Chain{Name: "test", Listener: dynamic.Listener{Depth: 1, Backoff: time.Second, ExternalTxEnabled: false}}
var config2 = &dynamic.Chain{Name: "test2", Listener: dynamic.Listener{Depth: 1, Backoff: time.Second, ExternalTxEnabled: false}}
var configExternTxEnabled = &dynamic.Chain{Name: "test2", Listener: dynamic.Listener{Depth: 1, Backoff: time.Second, ExternalTxEnabled: true}}

func TestCompareConfiguation(t *testing.T) {

	testSample := []struct {
		oldConfig       *dynamic.Configuration
		newConfig       *dynamic.Configuration
		expectedCommand []*Command
	}{
		{
			oldConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			expectedCommand: []*Command{},
		},
		{
			oldConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test":  config1,
				"test2": config1,
			}},
			expectedCommand: []*Command{
				{
					Type:  START,
					Chain: config1,
				},
			},
		},
		{
			oldConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test":  config1,
				"test2": config1,
			}},
			newConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			expectedCommand: []*Command{
				{
					Type:  STOP,
					Chain: config1,
				},
			},
		},
		{
			oldConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config2,
			}},
			expectedCommand: []*Command{
				{
					Type:  UPDATE,
					Chain: config2,
				},
			},
		},
		{
			oldConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": configExternTxEnabled,
			}},
			expectedCommand: []*Command{
				{
					Type:  UPDATE,
					Chain: configExternTxEnabled,
				},
			},
		},
	}

	for i, test := range testSample {
		command := CompareConfiguration(test.oldConfig, test.newConfig)
		assert.True(t, reflect.DeepEqual(command, test.expectedCommand), "%d/%d - should get the same command list - got %v - but expected %v", i+1, len(testSample), command, test.expectedCommand)
	}
}
