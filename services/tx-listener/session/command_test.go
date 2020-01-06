package session

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

var config1 = &dynamic.Node{Name: "test", Listener: &dynamic.Listener{Depth: 1, Backoff: time.Second}}
var config2 = &dynamic.Node{Name: "test2", Listener: &dynamic.Listener{Depth: 1, Backoff: time.Second}}

func TestCompareConfiguation(t *testing.T) {

	testSample := []struct {
		oldConfig       *dynamic.Configuration
		newConfig       *dynamic.Configuration
		expectedCommand []*Command
	}{
		{
			oldConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config1,
			}},
			expectedCommand: []*Command{},
		},
		{
			oldConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test":  config1,
				"test2": config1,
			}},
			expectedCommand: []*Command{
				{
					Type: START,
					Node: config1,
				},
			},
		},
		{
			oldConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test":  config1,
				"test2": config1,
			}},
			newConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config1,
			}},
			expectedCommand: []*Command{
				{
					Type: STOP,
					Node: config1,
				},
			},
		},
		{
			oldConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config1,
			}},
			newConfig: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": config2,
			}},
			expectedCommand: []*Command{
				{
					Type: UPDATE,
					Node: config2,
				},
			},
		},
	}

	for i, test := range testSample {
		command := CompareConfiguration(test.oldConfig, test.newConfig)
		assert.True(t, reflect.DeepEqual(command, test.expectedCommand), "%d/%d - should get the same command list - got %v - but expected %v", i+1, len(testSample), command, test.expectedCommand)
	}
}
