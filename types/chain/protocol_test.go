package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	p := Protocol{Name: "pantheon"}
	assert.Equal(t, true, p.IsPantheon(), "IsPantheon should be true")

	p = Protocol{Name: "quorum", Tag: "1.1.18"}
	assert.Equal(t, true, p.IsQuorum(), "IsQuorum should be true")

	p = Protocol{Name: "quorum", Tag: "2.2.2-alpha.1"}
	isTessera, err := p.IsTessera()
	assert.Nil(t, err, "IsTessera should not error")
	assert.Equal(t, true, isTessera, "IsTessera should be true")
}
