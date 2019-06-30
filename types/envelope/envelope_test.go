package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

func TestEnvelope(t *testing.T) {
	envelope := &Envelope{}
	assert.Equal(t, "", envelope.Error(), "Error message should be correct")

	envelope = &Envelope{
		Errors: []*err.Error{
			&err.Error{Code: 1, Message: "Timeout error", Component: "foo"},
			&err.Error{Code: 0, Message: "Unknown error", Component: "bar"},
		},
	}
	assert.Equal(t, `["00001@foo: Timeout error" "00000@bar: Unknown error"]`, envelope.Error(), "Error message should be correct")
}
