package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/error"
)

func TestEnvelope(t *testing.T) {
	envelope := &Envelope{}
	assert.Equal(t, "", envelope.Error(), "Error message should be correct")

	envelope = &Envelope{
		Errors: []*error.Error{
			{Code: 1, Message: "Timeout error", Component: "foo"},
			{Code: 0, Message: "Unknown error", Component: "bar"},
		},
	}
	assert.Equal(t, `["00001@foo: Timeout error" "00000@bar: Unknown error"]`, envelope.Error(), "Error message should be correct")

	// Test set and retrieving envelope Metadata
	_, ok := envelope.GetMetadataValue("test-key")
	assert.False(t, ok, "when no metadata has been set GetMetadataValue should not find data")

	envelope.SetMetadataValue("test-key", "test-value")
	v, ok := envelope.GetMetadataValue("test-key")
	assert.True(t, ok, "when metadata has been set GetMetadataValue should find data")
	assert.Equal(t, "test-value", v, "when metadata has been set GetMetadataValue should return expected value")
}
