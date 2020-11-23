package httpcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/httpcache/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

var generatedKey = "generatedKey"
var keySuffix = "keySuffix"

func TestHTTPCache_SetCacheValueSuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mockhandler.NewMockHandler(ctrl)
	cManager := mocks.NewMockCacheManager(ctrl)

	httpCache := newHTTPCache(cManager, testCacheRequest, testCacheResponse, keySuffix)
	h := httpCache.Handler(mockHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)

	cacheKey := fmt.Sprintf("%s-%s", keySuffix, generatedKey)
	cManager.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, false)
	cManager.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any())
	cManager.EXPECT().TTL()
	mockHandler.EXPECT().ServeHTTP(gomock.AssignableToTypeOf(&httptest.ResponseRecorder{}), req)

	h.ServeHTTP(rw, req)

	assert.Equal(t, rw.Result().StatusCode, 200)
}

func TestHTTPCache_SetCacheValueOnlyOnceOnConcurrentCalls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mockhandler.NewMockHandler(ctrl)
	cManager := mocks.NewMockCacheManager(ctrl)

	httpCache := newHTTPCache(cManager, testCacheRequest, testCacheResponse, keySuffix)
	h := httpCache.Handler(mockHandler)

	rw := httptest.NewRecorder()
	rw2 := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)
	res := newResponse([]byte(`responseBody`), http.Header{}, 200)
	bres, _ := res.toBytes()

	cacheKey := fmt.Sprintf("%s-%s", keySuffix, generatedKey)
	gomock.InOrder(
		cManager.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, false),
		cManager.EXPECT().Get(gomock.Any(), cacheKey).Return(bres, true),
	)
	cManager.EXPECT().TTL().Times(2)
	mockHandler.EXPECT().ServeHTTP(gomock.AssignableToTypeOf(&httptest.ResponseRecorder{}), req).Times(1)
	cManager.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any()).Times(1)

	utils.InParallel(
		func() { h.ServeHTTP(rw, req) },
		func() { h.ServeHTTP(rw2, req) },
	)

	assert.Equal(t, rw.Result().StatusCode, 200)
	assert.Equal(t, rw2.Result().StatusCode, 200)
	assert.NotEqual(t, rw2.Result().Header["X-Cache-Control"], rw.Result().Header["X-Cache-Control"])
}

func TestHTTPCache_GetCacheValueSuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mockhandler.NewMockHandler(ctrl)
	cManager := mocks.NewMockCacheManager(ctrl)

	httpCache := newHTTPCache(cManager, testCacheRequest, testCacheResponse, keySuffix)
	h := httpCache.Handler(mockHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)

	res := newResponse([]byte(`responseBody`), http.Header{}, 200)
	bres, _ := res.toBytes()
	cacheKey := fmt.Sprintf("%s-%s", keySuffix, generatedKey)
	cManager.EXPECT().TTL()
	cManager.EXPECT().Get(gomock.Any(), cacheKey).Return(bres, true)

	h.ServeHTTP(rw, req)

	result := rw.Result()
	defer result.Body.Close()
	rBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, result.StatusCode, 200)
	assert.NotEmpty(t, result.Header.Get("X-Cache-Control"))
	assert.Equal(t, rBody, []byte("responseBody"))
}

func testCacheRequest(_ *http.Request) (isCached bool, key string, ttl time.Duration, err error) {
	return true, generatedKey, 0, nil
}

func testCacheResponse(_ *http.Response) bool {
	return true
}
