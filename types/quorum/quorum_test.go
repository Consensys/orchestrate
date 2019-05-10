package quorum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuorum(t *testing.T) {
	q := &Quorum{Version: "2.2.3-alpha.1"}
	is, err := q.IsTessera()
	assert.NoError(t, err, "#1 should not error")
	assert.True(t, is, "#1: Should be Tessera compatible")

	q = &Quorum{Version: "2.2.2-alpha.1"}
	is, err = q.IsTessera()
	assert.NoError(t, err, "#2 should not error")
	assert.True(t, is, "#2: Should not be Tessera compatible")
}
