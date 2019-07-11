package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocol(t *testing.T) {
	p := Protocol{Type: ProtocolType_PANTHEON}
	assert.Equal(t, true, p.IsPantheon(), "IsPantheon should be true")

	p = Protocol{Type: ProtocolType_QUORUM, Tag: "1.1.18"}
	assert.Equal(t, true, p.IsQuorum(), "IsQuorum should be true")

	p = Protocol{Type: ProtocolType_QUORUM, Tag: "2.2.2-alpha.1"}
	isTessera, err := p.IsTessera()
	assert.Nil(t, err, "IsTessera should not error")
	assert.Equal(t, true, isTessera, "IsTessera should be true")
}
