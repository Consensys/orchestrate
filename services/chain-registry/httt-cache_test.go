package chainregistry

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
)

func TestHTTPCacheRequest_Valid(t *testing.T) {
	msg := rpc.JSONRpcMessage{
		Method: "eth_getBlockByNumber",
	}
	msg.Params, _ = json.Marshal([]string{"0x0", "false"})

	body, _ := json.Marshal(msg)
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(body))

	c, k, err := httpCacheRequest(req)
	assert.NoError(t, err)
	assert.True(t, c)
	assert.Equal(t, "_-eth_getBlockByNumber([\"0x0\",\"false\"])", k)
}

func TestHTTPCacheRequest_IgnoreReqType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	c, _, err := httpCacheRequest(req)
	assert.NoError(t, err)
	assert.False(t, c)
}

func TestHTTPCacheRequest_IgnoreRPCMethod(t *testing.T) {
	msg := rpc.JSONRpcMessage{
		Method: "eth_getTransactionCount",
	}

	body, _ := json.Marshal(msg)
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(body))

	c, _, err := httpCacheRequest(req)
	assert.NoError(t, err)
	assert.False(t, c)
}
