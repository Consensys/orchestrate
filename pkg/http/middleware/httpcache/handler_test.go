package httpcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/httpcache/mocks"
)

var generatedKey = "generatedKey"
var keySuffix = "keySuffix"

func TestSetCacheFlow_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mockhandler.NewMockHandler(ctrl)
	cManager := mocks.NewMockCacheManager(ctrl)

	httpCache := newHTTPCache(cManager, generateKey, keySuffix)
	h := httpCache.Handler(mockHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)

	cacheKey := fmt.Sprintf("%s-%s", generatedKey, keySuffix)
	cManager.EXPECT().Get(gomock.Any(), cacheKey).Return(nil, false)
	cManager.EXPECT().Set(gomock.Any(), cacheKey, gomock.Any())
	cManager.EXPECT().TTL()
	mockHandler.EXPECT().ServeHTTP(gomock.AssignableToTypeOf(&httptest.ResponseRecorder{}), req)

	h.ServeHTTP(rw, req)

	assert.Equal(t, rw.Result().StatusCode, 200)
}

func TestGetCacheFlow_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mockhandler.NewMockHandler(ctrl)
	cManager := mocks.NewMockCacheManager(ctrl)

	httpCache := newHTTPCache(cManager, generateKey, keySuffix)
	h := httpCache.Handler(mockHandler)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://host.com/live", nil)

	res := newResponse([]byte(`responseBody`), http.Header{}, 200)
	bres, _ := res.toBytes()
	cacheKey := fmt.Sprintf("%s-%s", generatedKey, keySuffix)
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

func generateKey(_ *http.Request) (isCached bool, key string, err error) {
	return true, generatedKey, nil
}
