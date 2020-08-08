package httpcache

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dgraph-io/ristretto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type Builder struct {
	cache   *ristretto.Cache
	cReqKey CacheRequest
}

type CacheRequest func(req *http.Request) (isCached bool, key string, err error)

func NewBuilder(cache *ristretto.Cache, fReq CacheRequest) *Builder {
	return &Builder{
		cache:   cache,
		cReqKey: fReq,
	}
}

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.HTTPCache)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	cManager := newManager(b.cache, cfg.TTL)

	m := newHTTPCache(cManager, b.cReqKey, cfg.KeySuffix)
	return m.Handler, nil, nil
}

type HTTPCache struct {
	cManager CacheManager
	cReqKey  CacheRequest
	cSuffix  string
}

func newHTTPCache(cManager CacheManager, cReqKey CacheRequest, cSuffix string) *HTTPCache {
	return &HTTPCache{
		cManager: cManager,
		cReqKey:  cReqKey,
		cSuffix:  cSuffix,
	}
}

func (cm *HTTPCache) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cacheActive, cacheKeyBase, err := cm.cacheRequest(req)
		if !cacheActive {
			if err != nil {
				log.WithContext(req.Context()).WithError(err).Error("HTTPCache: errors were found")
			}

			log.WithContext(req.Context()).Debugf("HTTPCache: request is skipped")
			h.ServeHTTP(rw, req)
			return
		}

		cacheKey := fmt.Sprintf("%s-%s", cacheKeyBase, cm.cSuffix)
		b, ok := cm.cManager.Get(req.Context(), cacheKey)
		if ok {
			res, err := bytesToResponse((b.([]byte)))
			if err != nil {
				log.WithContext(req.Context()).WithError(err).Error("HTTPCache: errors were found")
			} else {
				for k, v := range res.Header {
					rw.Header().Set(k, strings.Join(v, ","))
				}
				rw.Header().Set("X-Cache-Control", fmt.Sprintf("max-age=%dms", cm.cManager.TTL().Milliseconds()))
				rw.WriteHeader(res.StatusCode)
				if _, err := rw.Write(res.Value); err != nil {
					log.WithContext(req.Context()).WithError(err).Error("HTTPCache: errors were found")
				}

				log.WithContext(req.Context()).WithField("key", cacheKey).
					Debugf("HTTPCache: response fetch from cache")

				return
			}
		}

		// Otherwise, capture response and cache it
		rwRecoder := httptest.NewRecorder()
		h.ServeHTTP(rwRecoder, req)

		// Extract response content a write it into response
		result := rwRecoder.Result()
		if result.StatusCode == 200 {
			r := newResponse(rwRecoder.Body.Bytes(), result.Header, result.StatusCode)
			b, err := r.toBytes()
			if err != nil {
				log.WithContext(req.Context()).WithError(err).Error("HTTPCache: errors were found")
				return
			}

			cm.cManager.Set(req.Context(), cacheKey, b)
		}

		for k, v := range result.Header {
			rw.Header().Set(k, strings.Join(v, ","))
		}

		rw.WriteHeader(result.StatusCode)
		if _, err := rw.Write(rwRecoder.Body.Bytes()); err != nil {
			log.WithContext(req.Context()).WithError(err).Error("HTTPCache: errors were found")
		}

		log.WithContext(req.Context()).WithField("key", cacheKey).
			Debugf("HTTPCache: caching response for %dms", cm.cManager.TTL().Milliseconds())
	})
}

func (cm *HTTPCache) cacheRequest(req *http.Request) (c bool, k string, err error) {
	if req.Header.Get("X-Cache-Control") == "no-cache" {
		return false, "", nil
	}

	return cm.cReqKey(req)
}
