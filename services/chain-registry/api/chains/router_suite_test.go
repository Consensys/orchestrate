package chains

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mocks"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

const (
	expectedInternalServerErrorBody  = "{\"message\":\"test error\"}\n"
	expectedNotFoundErrorBody        = "{\"message\":\"DB200@: not found error\"}\n"
	expectedInvalidBodyError         = "{\"message\":\"FF000@chain-registry.store.api: invalid body\"}\n"
	expectedUnknownBodyError         = "{\"message\":\"FF000@chain-registry.store.api: json: unknown field \\\"unknownField\\\"\"}\n"
	expectedSuccessStatusBody        = "{}\n"
	expectedSuccessStatusSliceBody   = "[]\n"
	expectedSuccessStatusContentType = "application/json"
	expectedErrorStatusContentType   = "text/plain; charset=utf-8"
	notFoundErrorFilter              = "notFoundError"
)

type HTTPRouteTests struct {
	name                string
	store               func(t *testing.T) models.ChainRegistryStore
	httpMethod          string
	path                string
	body                func() []byte
	expectedStatusCode  int
	expectedContentType string
	expectedBody        func() string
}

func TestHTTPRouteTests(t *testing.T) {
	t.Parallel()

	testsSuite := [][]HTTPRouteTests{
		deleteChainByUUIDTests,
		getChainsTests,
		getChainsByUUIDTests,
		patchChainByUUIDTests,
		postChainTests,
	}

	for _, tests := range testsSuite {
		for _, tt := range tests {
			tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel() // marks each test case as capable of running in parallel with each other

				router := mux.NewRouter()
				NewHandler(tt.store(t)).Append(router)

				// Normal tests
				writer := httptest.NewRecorder()
				request := httptest.NewRequest(tt.httpMethod, tt.path, bytes.NewReader(tt.body()))

				router.ServeHTTP(writer, request)
				testResponse(t, writer, tt.expectedStatusCode, tt.expectedContentType, tt.expectedBody())
			})
		}
	}

	// Multi-tenant tests
	for _, tests := range testsSuite {
		for _, tt := range tests {
			tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel() // marks each test case as capable of running in parallel with each other

				router := mux.NewRouter()
				NewHandler(tt.store(t)).Append(router)

				writer := httptest.NewRecorder()
				request := httptest.NewRequest(tt.httpMethod, tt.path, bytes.NewReader(tt.body()))
				request = request.WithContext(multitenancy.WithTenantID(request.Context(), "tenantID1"))
				router.ServeHTTP(writer, request)
				testResponse(t, writer, tt.expectedStatusCode, tt.expectedContentType, tt.expectedBody())
			})
		}
	}

}

func testResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int, expectedContentType, expectedBody string) {
	assert.Equal(t, expectedContentType, w.Header().Get("Content-Type"), "Did not get expected content type %s, but got %s", expectedContentType, w.Header().Get("Content-Type"))
	assert.Equal(t, expectedStatusCode, w.Code, "Did not get expected HTTP status code %d, but got %d", expectedStatusCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String(), "Did not get expected body %s, but got %s", expectedBody, w.Body.String())
}

func UseMockChainRegistry(t *testing.T) models.ChainRegistryStore {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mocks.NewMockChainRegistryStore(mockCtrl)

	mockStore.EXPECT().RegisterChain(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, chain *models.Chain) error {
			listenerDepth := uint64(1)
			listenerBlockPosition := int64(1)
			listenerFromBlock := int64(1)
			listenerBackOffDuration := "1s"
			chain.UUID = "1"
			chain.Name = "chainName1"
			chain.TenantID = multitenancy.TenantIDFromContext(ctx)
			chain.URLs = []string{"testUrl1", "testUrl2"}
			chain.ListenerDepth = &listenerDepth
			chain.ListenerBlockPosition = &listenerBlockPosition
			chain.ListenerFromBlock = &listenerFromBlock
			chain.ListenerBackOffDuration = &listenerBackOffDuration
			return nil
		}).AnyTimes()

	mockStore.EXPECT().GetChains(gomock.Any(), gomock.Any()).Return([]*models.Chain{}, nil).AnyTimes()
	mockStore.EXPECT().GetChainsByTenant(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*models.Chain{}, nil).AnyTimes()

	mockStore.EXPECT().GetChainByUUID(gomock.Any(), gomock.Any()).Return(&models.Chain{}, nil).AnyTimes()
	mockStore.EXPECT().GetChainByUUIDAndTenant(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Chain{}, nil).AnyTimes()

	mockStore.EXPECT().UpdateChainByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().UpdateChainByUUID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockStore.EXPECT().DeleteChainByUUID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().DeleteChainByUUIDAndTenant(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	return mockStore
}

var errTest = fmt.Errorf("test error")
var errNotFound = errors.NotFoundError("not found error")

func UseErrorChainRegistry(t *testing.T) models.ChainRegistryStore {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mocks.NewMockChainRegistryStore(mockCtrl)

	mockStore.EXPECT().RegisterChain(gomock.Any(), gomock.Any()).Return(errTest).AnyTimes()

	mockStore.EXPECT().GetChains(gomock.Any(), gomock.Any()).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().GetChainsByTenant(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errTest).AnyTimes()

	mockStore.EXPECT().GetChainByUUID(gomock.Any(), gomock.Eq("0")).Return(nil, errNotFound).AnyTimes()
	mockStore.EXPECT().GetChainByUUIDAndTenant(gomock.Any(), gomock.Eq("0"), gomock.Any()).Return(nil, errNotFound).AnyTimes()

	mockStore.EXPECT().GetChainByUUID(gomock.Any(), gomock.Not(gomock.Eq("0"))).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().GetChainByUUIDAndTenant(gomock.Any(), gomock.Not(gomock.Eq("0")), gomock.Any()).Return(nil, errTest).AnyTimes()

	mockStore.EXPECT().UpdateChainByName(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, chain *models.Chain) error {
			if chain.Name == notFoundErrorFilter {
				return errNotFound
			}
			return errTest
		}).AnyTimes()
	mockStore.EXPECT().UpdateChainByUUID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, chain *models.Chain) error {
			if chain.UUID == "0" {
				return errNotFound
			}
			return errTest
		}).AnyTimes()

	mockStore.EXPECT().DeleteChainByUUID(gomock.Any(), gomock.Eq("0")).Return(errNotFound).AnyTimes()
	mockStore.EXPECT().DeleteChainByUUIDAndTenant(gomock.Any(), gomock.Eq("0"), gomock.Any()).Return(errNotFound).AnyTimes()

	mockStore.EXPECT().DeleteChainByUUID(gomock.Any(), gomock.Not(gomock.Eq("0"))).Return(errTest).AnyTimes()
	mockStore.EXPECT().DeleteChainByUUIDAndTenant(gomock.Any(), gomock.Not(gomock.Eq("0")), gomock.Any()).Return(errTest).AnyTimes()

	return mockStore
}
