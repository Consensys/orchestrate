package proxy

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/errors"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	"github.com/consensys/orchestrate/pkg/utils"
)

var rpcCachedMethods = map[string]bool{
	"eth_getBlockByNumber":      true,
	"eth_getTransactionReceipt": true,
}

func HTTPCacheRequest(ctx context.Context, req *http.Request) (c bool, k string, ttl time.Duration, err error) {
	logger := log.FromContext(ctx)
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
		logger.Debugf("request is not an rpc")
		return false, "", 0, nil
	}

	if _, ok := rpcCachedMethods[msg.Method]; !ok {
		logger.WithField("method", msg.Method).Debug("rpc method is ignored")
		return false, "", 0, nil
	}

	cacheKey := fmt.Sprintf("%s(%s)", msg.Method, string(msg.Params))
	if msg.Method == "eth_getBlockByNumber" && strings.Contains(string(msg.Params), "latest") {
		return true, cacheKey, time.Second, nil
	}

	return true, cacheKey, 0, nil
}

func HTTPCacheResponse(ctx context.Context, resp *http.Response) bool {
	logger := log.FromContext(ctx)
	var msg ethclient.JSONRpcMessage
	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	err := json.UnmarshalBody(reader, &msg)
	if err != nil {
		logger.WithError(err).Debug("failed to decode rpc response")
		return false
	}

	if msg.Error != nil {
		logger.WithField("error", msg.Error.Message).Debug("skipped rpc error responses")
		return false
	}

	if len(msg.Result) == 0 {
		logger.Debug("skipped rpc empty response results")
		return false
	}

	return true
}

func httpCacheGenerateChainKey(chain *entities.Chain) string {
	// Order urls to identify common chain definitions
	sort.Sort(utils.Alphabetic(chain.URLs))

	var urls = strings.Join(chain.URLs, "_")

	if chain.PrivateTxManager != nil {
		urls = urls + "_" + chain.PrivateTxManager.URL
	}

	return urls
}
