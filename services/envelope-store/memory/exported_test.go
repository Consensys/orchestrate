// +build unit

package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, GlobalEnvelopeStore(), "Global should have been set")

	var envelopeStore *EnvelopeStore
	SetGlobalEnvelopeStore(envelopeStore)
	assert.Nil(t, GlobalEnvelopeStore(), "Global should be reset to nil")
}
