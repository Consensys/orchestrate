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
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type CacheRequest func(req *http.Request) (isCached bool, key string, ttl time.Duration, err error)
type CacheResponse func(res *http.Response) bool

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

	m := newHTTPCache(cManager, b.cacheReq, b.cacheRes, cfg.KeySuffix)
	return m.Handler, nil, nil
}

type HTTPCache struct {
	cManager CacheManager
	cacheReq CacheRequest
	cacheRes CacheResponse
	cSuffix  string
	reqMutex map[uint8]*sync.Mutex
	mutex    *sync.RWMutex
}

func newHTTPCache(cManager CacheManager, cacheReq CacheRequest, cacheRes CacheResponse, cSuffix string) *HTTPCache {
	return &HTTPCache{
		cManager: cManager,
		cacheReq: cacheReq,
		cacheRes: cacheRes,
		cSuffix:  cSuffix,
		mutex:    &sync.RWMutex{},
		reqMutex: make(map[uint8]*sync.Mutex),
	}
}

func (cm *HTTPCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger := log.WithContext(req.Context())
		cacheActive, cacheKeyBase, ttl, err := cm.cacheRequest(req)
		if err != nil {
			logger.WithError(err).Error("HTTPCache: errors were found")
		}

		if !cacheActive {
			logger.Debugf("HTTPCache: request is skipped")
			h.ServeHTTP(rw, req)
			return
		}

		cacheKey := fmt.Sprintf("%s-%s", cm.cSuffix, cacheKeyBase)
		logger = logger.WithField("key", cacheKey)

		cMutex := cm.distributedRequestMutex(cacheKey)
		cMutex.Lock()
		defer cMutex.Unlock()

		b, ok := cm.cManager.Get(req.Context(), cacheKey)
		if ok {
			res, err := bytesToResponse(b)
			if err != nil {
				logger.WithError(err).Error("HTTPCache: errors were found")
			} else {
				for k, v := range res.Header {
					rw.Header().Set(k, strings.Join(v, ","))
				}
				rw.Header().Set("X-Cache-Control", fmt.Sprintf("max-age=%dms", cm.cManager.TTL().Milliseconds()))
				rw.WriteHeader(res.StatusCode)
				if _, err := rw.Write(res.Value); err != nil {
					logger.WithError(err).Error("HTTPCache: errors were found")
				}

				logger.Debugf("HTTPCache: response fetched from cache")
				return
			}
		}

		// Otherwise, capture response and cache it
		rwRecoder := httptest.NewRecorder()
		h.ServeHTTP(rwRecoder, req)

		// Extract response content a write it into response, only successful responses are cached
		result := rwRecoder.Result()
		if cm.cacheResponse(result) {
			r := newResponse(rwRecoder.Body.Bytes(), result.Header, result.StatusCode)
			b, err := r.toBytes()
			if err != nil {
				logger.WithError(err).Error("HTTPCache: errors were found")
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
			logger.WithError(err).Error("HTTPCache: errors were found")
		}
	})
}

func (cm *HTTPCache) cacheRequest(req *http.Request) (c bool, k string, ttl time.Duration, err error) {
	if req.Header.Get("X-Cache-Control") == "no-cache" {
		return false, "", 0, nil
	}

	return cm.cacheReq(req)
}

func (cm *HTTPCache) cacheResponse(res *http.Response) bool {
	if res.StatusCode != 200 {
		log.WithField("status", res.StatusCode).Debugf("HTTPCache: skip not 200 responses")
		return false
	}

	return cm.cacheRes(res)
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
