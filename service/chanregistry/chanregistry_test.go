package chanregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChanRegistry(t *testing.T) {
	chanregistry := NewChanRegistry()
	assert.NotNil(t, chanregistry, "Registry should have been set")

	scenarioID := "testScenario"
	topic := "testTopic"
	topicChan := chanregistry.NewEnvelopeChan(scenarioID, topic)
	assert.NotNil(t, topicChan, "Channel should have been created")

	assert.Nil(t, chanregistry.GetEnvelopeChan("dummy", "dummy"), "Should not get a chan")
	assert.NotNil(t, chanregistry.GetEnvelopeChan(scenarioID, topic), "Should get a chan")

	_ = chanregistry.CloseEnvelopeChan(scenarioID, topic)
	assert.Nil(t, chanregistry.GetEnvelopeChan(scenarioID, topic), "Should not get a chan")

}
