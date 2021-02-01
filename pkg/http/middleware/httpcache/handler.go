package httpcache

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
)

const component = "http.cache"

type CacheRequest func(ctx context.Context, req *http.Request) (isCached bool, key string, ttl time.Duration, err error)
type CacheResponse func(ctx context.Context, res *http.Response) bool

type Builder struct {
	cache    *ristretto.Cache
	cacheReq CacheRequest
	cacheRes CacheResponse
}

func NewBuilder(cache *ristretto.Cache, cacheReq CacheRequest, cacheRes CacheResponse) *Builder {
	return &Builder{
		cache:    cache,
		cacheReq: cacheReq,
		cacheRes: cacheRes,
	}
}

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.HTTPCache)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	cManager := newManager(b.cache, cfg.TTL)
	logger := log.NewLogger().SetComponent(component)
	logger.Debug("middleware built successfully")

	m := newHTTPCache(cManager, b.cacheReq, b.cacheRes, cfg.KeySuffix, logger)
	return m.Handler, nil, nil
}

type HTTPCache struct {
	cManager CacheManager
	cacheReq CacheRequest
	cacheRes CacheResponse
	cSuffix  string
	reqMutex map[uint8]*sync.Mutex
	mutex    *sync.RWMutex
	logger   *log.Logger
}

func newHTTPCache(cManager CacheManager, cacheReq CacheRequest, cacheRes CacheResponse, cSuffix string, logger *log.Logger) *HTTPCache {
	return &HTTPCache{
		cManager: cManager,
		cacheReq: cacheReq,
		cacheRes: cacheRes,
		cSuffix:  cSuffix,
		mutex:    &sync.RWMutex{},
		reqMutex: make(map[uint8]*sync.Mutex),
		logger:   logger,
	}
}

func (cm *HTTPCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger := cm.logger.WithContext(req.Context())
		ctx := log.With(req.Context(), logger)

		cacheActive, cacheKeyBase, ttl, err := cm.cacheRequest(ctx, req)
		if err != nil {
			logger.WithError(err).Error("failed to build cache request")
		}

		if !cacheActive {
			h.ServeHTTP(rw, req)
			return
		}

		cacheKey := fmt.Sprintf("%s-%s", cm.cSuffix, cacheKeyBase)
		logger = logger.WithField("key", cacheKey)

		cMutex := cm.distributedRequestMutex(cacheKey)
		cMutex.Lock()
		defer cMutex.Unlock()

		b, ok := cm.cManager.Get(ctx, cacheKey)
		if ok {
			res, err := bytesToResponse(b)
			if err != nil {
				logger.WithError(err).Error("failed to decode response")
			} else {
				for k, v := range res.Header {
					rw.Header().Set(k, strings.Join(v, ","))
				}
				rw.Header().Set("X-Cache-Control", fmt.Sprintf("max-age=%dms", cm.cManager.TTL().Milliseconds()))
				rw.WriteHeader(res.StatusCode)
				if _, err := rw.Write(res.Value); err != nil {
					logger.WithError(err).Error("failed to write cache")
				}

				logger.Debug("response was pull from cache")
				return
			}
		}

		// Otherwise, capture response and cache it
		rwRecoder := httptest.NewRecorder()
		h.ServeHTTP(rwRecoder, req)

		// Extract response content a write it into response, only successful responses are cached
		result := rwRecoder.Result()
		if cm.cacheResponse(ctx, result) {
			r := newResponse(rwRecoder.Body.Bytes(), result.Header, result.StatusCode)
			b, err := r.toBytes()
			if err != nil {
				logger.WithError(err).Error("failed to write cached response")
				return
			}

			// In case we have a customize TTL
			if ttl != 0 {
				logger = logger.WithField("ttl", ttl.String())
				cm.cManager.SetWithTTL(req.Context(), cacheKey, b, ttl)
			} else {
				logger = logger.WithField("ttl", cm.cManager.TTL().String())
				cm.cManager.Set(req.Context(), cacheKey, b)
			}

		}

		for k, v := range result.Header {
			rw.Header().Set(k, strings.Join(v, ","))
		}

		rw.WriteHeader(result.StatusCode)
		if _, err := rw.Write(rwRecoder.Body.Bytes()); err != nil {
			logger.WithError(err).Error("failed to write response")
		}
	})
}

func (cm *HTTPCache) cacheRequest(ctx context.Context, req *http.Request) (c bool, k string, ttl time.Duration, err error) {
	if req.Header.Get("X-Cache-Control") == "no-cache" {
		return false, "", 0, nil
	}

	return cm.cacheReq(ctx, req)
}

func (cm *HTTPCache) cacheResponse(ctx context.Context, res *http.Response) bool {
	if res.StatusCode != 200 {
		cm.logger.WithField("status", res.StatusCode).Debug("skip responses with status code different than 200")
		return false
	}

	return cm.cacheRes(ctx, res)
}

// Generate/Retrieve mutex item to synchronize the access to cached request/responses
func (cm *HTTPCache) distributedRequestMutex(cacheKey string) *sync.Mutex {
	mutexHashKey := generateCacheMutexKey(cacheKey)
	cm.mutex.RLock()
	cMutex := cm.reqMutex[mutexHashKey]
	cm.mutex.RUnlock()
	if cMutex == nil {
		// In case targeted mutex key remains nil after exclusive access, we update it
		cm.mutex.Lock()
		if cm.reqMutex[mutexHashKey] == nil {
			cMutex = &sync.Mutex{}
			cm.reqMutex[mutexHashKey] = cMutex
		}
		cm.mutex.Unlock()
	}

	return cMutex
}

// Simple hash distribution function of cacheKey values
// Source: https://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
func generateCacheMutexKey(cacheKey string) uint8 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(cacheKey))
	sum := h.Sum32()
	return uint8(sum % 256)
}
