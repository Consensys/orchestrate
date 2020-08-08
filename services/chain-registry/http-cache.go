package chainregistry

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

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
	"eth_getBlockNumber":        true,
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
	}

	cacheKey, err := httpCacheGenerateKey(req.Context(), &msg)
	if err != nil {
		log.WithContext(req.Context()).WithError(err).Debugf("HTTPCache: cannot generate cache key")
		return false, "", nil
	}

	return true, cacheKey, nil
}

func httpCacheGenerateKey(ctx context.Context, msg *rpc.JSONRpcMessage) (string, error) {
	jsonParams, err := json.Marshal(msg.Params)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(string(jsonParams) + msg.Method))

	key := fmt.Sprintf("%s%s", multitenancy.TenantIDFromContext(ctx), hex.EncodeToString(hash[:]))
	return key, nil
}

func httpCacheGenerateChainKey(chain *models.Chain) string {
	sort.Sort(utils.Alphabetic(chain.URLs))
	urls, err := json.Marshal(chain.URLs)
	// We do not intend to fail in case of an marshaling issue, instead we use ChainID as backup key
	if err != nil {
		return chain.ChainID
	}

	hash := md5.Sum(urls)
	return hex.EncodeToString(hash[:])
}
