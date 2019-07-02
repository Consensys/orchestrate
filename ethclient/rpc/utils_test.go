package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

func TestParsesPantheonClientVersion(t *testing.T) {
	clientType := ClientTypeParser("pantheon/1.1.1")
	assert.Equal(t, types.PantheonClient, clientType)
}

func TestParsesGethClientVersion(t *testing.T) {
	clientType := ClientTypeParser("Geth/v1.8.18-stable-ef256cb2(quorum-v2.2.3)/linux-amd64/go1.10.1")
	assert.Equal(t, types.QuorumClient, clientType, "%s != %s", types.QuorumClient, clientType)
}

func TestParseUnknownClientVersion(t *testing.T) {
	clientType := ClientTypeParser("unknown-client")
	assert.Equal(t, types.UnknownClient, clientType)
}

var clientDetectionTests = []struct {
	clientVersion string
	clientType    types.ClientType
}{
	{"pantheon/1.1.1", types.PantheonClient},
	{"Geth/v1.8.18-stable-ef256cb2(quorum-v2.2.3)/linux-amd64/go1.10.1", types.QuorumClient},
	{"unknown-client", types.UnknownClient},
}

func TestClientDetection(t *testing.T) {
	for _, tt := range clientDetectionTests {
		t.Run(tt.clientVersion, func(t *testing.T) {
			detectedClientType := ClientTypeParser(tt.clientVersion)
			assert.Equal(t, tt.clientType, detectedClientType, "%s != %s", types.QuorumClient, tt.clientType)
		})
	}
}
