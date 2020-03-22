// +build unit

package sarama

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	err "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
)

func TestNewClient(t *testing.T) {
	_, e := NewClient([]string{"unknown"}, sarama.NewConfig())
	assert.Error(t, e, "Client should error")
	ie, ok := e.(*err.Error)
	assert.True(t, ok, "Error should cast to internal error")
	assert.Equal(t, "broker.sarama", ie.GetComponent(), "Component should be correct")
	assert.True(t, errors.IsConnectionError(ie), "Error should be a connection error")
}
