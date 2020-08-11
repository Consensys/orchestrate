package chainregistry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/utils"
)

var rpcCachedMethods = map[string]bool{
	"eth_getBlockByNumber":      true,
	"eth_getTransactionReceipt": true,
}

func httpCacheRequest(req *http.Request) (c bool, k string, err error) {
	if req.Method != "POST" {
		return false, "", nil
	}

	if req.Body == nil {
		return false, "", nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, "", errors.InternalError("can't read request body: %q", err)
	}

	// And now set a new body, which will simulate the same data we read
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var msg rpc.JSONRpcMessage
	err = json.Unmarshal(body, &msg)
	// In case request does not correspond to one of expected call RPC call, we ignore
	if err != nil {
		log.WithContext(req.Context()).Debugf("HTTPCache: request is not an RPC message")
		return false, "", nil
	}

	if _, ok := rpcCachedMethods[msg.Method]; !ok {
		log.WithContext(req.Context()).Debugf("HTTPCache: RPC method is ignored: %s", msg.Method)
		return false, "", nil
	} else if msg.Method == "eth_getBlockByNumber" && strings.Contains(string(msg.Params), "latest") {
		log.WithContext(req.Context()).Debugf("HTTPCache: call to the 'latest' block is ignored")
		return false, "", nil
	}

	cacheKey := httpCacheGenerateKey(req.Context(), &msg)

	return true, cacheKey, nil
}

func httpCacheResponse(req *http.Response) bool {
	var msg rpc.JSONRpcMessage
	err := json.UnmarshalBody(req.Body, &msg)
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

func httpCacheGenerateKey(ctx context.Context, msg *rpc.JSONRpcMessage) string {
	return fmt.Sprintf("%s-%s(%s)", multitenancy.TenantIDFromContext(ctx), msg.Method, string(msg.Params))
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
