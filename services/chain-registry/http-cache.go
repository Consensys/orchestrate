package chainregistry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/utils"
)

var rpcCachedMethods = map[string]bool{
	"eth_getBlockByNumber":      true,
	"eth_getTransactionReceipt": true,
}

func httpCacheRequest(req *http.Request) (c bool, k string, ttl time.Duration, err error) {
	logger := log.WithContext(req.Context())
	if req.Method != "POST" {
		return false, "", 0, nil
	}

	if req.Body == nil {
		return false, "", 0, nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, "", 0, errors.InternalError("can't read request body: %q", err)
	}

	// And now set a new body, which will simulate the same data we read
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var msg ethclient.JSONRpcMessage
	err = json.Unmarshal(body, &msg)
	// In case request does not correspond to one of expected call RPC call, we ignore
	if err != nil {
		logger.Debugf("HTTPCache: request is not an RPC message")
		return false, "", 0, nil
	}

	if _, ok := rpcCachedMethods[msg.Method]; !ok {
		logger.Debugf("HTTPCache: RPC method is ignored: %s", msg.Method)
		return false, "", 0, nil
	}

	cacheKey := httpCacheGenerateKey(req.Context(), &msg)

	if msg.Method == "eth_getBlockByNumber" && strings.Contains(string(msg.Params), "latest") {
		return true, cacheKey, time.Second, nil
	}

	return true, cacheKey, 0, nil
}

func httpCacheResponse(resp *http.Response) bool {
	var msg ethclient.JSONRpcMessage
	err := json.UnmarshalBody(resp.Body, &msg)
	if err != nil {
		log.WithError(err).Debugf("HTTPCache: cannot decode response")
		return false
	}

	if msg.Error != nil {
		log.WithField("error", msg.Error.Message).Debugf("HTTPCache: skip RPC error responses")
		return false
	}

	if len(msg.Result) == 0 {
		log.Debugf("HTTPCache: skip RPC empty response results")
		return false
	}

	return true
}

func httpCacheGenerateKey(_ context.Context, msg *ethclient.JSONRpcMessage) string {
	return fmt.Sprintf("%s(%s)", msg.Method, string(msg.Params))
}

func httpCacheGenerateChainKey(chain *models.Chain) string {
	// Order urls to identify common chain definitions
	sort.Sort(utils.Alphabetic(chain.URLs))

	var urls = strings.Join(chain.URLs, "_")

	var privURLs []string
	for _, p := range chain.PrivateTxManagers {
		privURLs = append(privURLs, p.URL)
	}

	if len(privURLs) > 0 {
		urls = urls + "_" + strings.Join(privURLs, "_")
	}

	return urls
}
