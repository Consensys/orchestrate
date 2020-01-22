package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

const (
	expectedInternalServerErrorBody  = "{\"message\":\"test error\"}\n"
	expectedNotFoundErrorBody        = "{\"message\":\"DB200@: not found error\"}\n"
	expectedInvalidURLErrorBody      = "{\"message\":\"FF000@chain-registry.store.api: parse test.com: invalid URI for request\"}\n"
	expectedInvalidErrorBody         = "{\"message\":\"FF000@chain-registry.store.api: json: unknown field \\\"unknownField\\\"\"}\n"
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
		deleteNodeByIDTests,
		deleteNodesByNameTests,
		getNodesTests,
		getNodesByIDTests,
		getNodesByNameTests,
		getNodesByTenantIDTests,
		patchNodeByIDTests,
		patchNodeByNameTests,
		patchBlockPositionByIDTests,
		patchBlockNumberByNameTests,
		postNodeTests,
	}

	for _, tests := range testsSuite {
		for _, tt := range tests {
			tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel() // marks each test case as capable of running in parallel with each other
				c := tt.store(t)
				r := mux.NewRouter()
				NewHandler(c).Append(r)

				w := httptest.NewRecorder()
				r.ServeHTTP(w, httptest.NewRequest(tt.httpMethod, tt.path, bytes.NewReader(tt.body())))

				testResponse(t, w, tt.expectedStatusCode, tt.expectedContentType, tt.expectedBody())
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

	mockStore.EXPECT().RegisterNode(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, node *models.Node) error {
			node.ID = "1"
			node.Name = "nodeName1"
			node.TenantID = "tenantID1"
			node.URLs = []string{"testUrl1", "testUrl2"}
			node.ListenerDepth = 1
			node.ListenerBlockPosition = 1
			node.ListenerFromBlock = 1
			node.ListenerBackOffDuration = "1s"
			return nil
		}).AnyTimes()
	mockStore.EXPECT().GetNodes(gomock.Any(), gomock.Any()).Return([]*models.Node{}, nil).AnyTimes()
	mockStore.EXPECT().GetNodesByTenantID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*models.Node{}, nil).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeName(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Node{}, nil).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Node{}, nil).AnyTimes()
	mockStore.EXPECT().GetNodeByID(gomock.Any(), gomock.Any()).Return(&models.Node{}, nil).AnyTimes()
	mockStore.EXPECT().UpdateNodeByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().UpdateNodeByID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().DeleteNodeByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockStore.EXPECT().DeleteNodeByID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	return mockStore
}

var errTest = fmt.Errorf("test error")
var errNotFound = errors.NotFoundError("not found error")

func UseErrorChainRegistry(t *testing.T) models.ChainRegistryStore {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mocks.NewMockChainRegistryStore(mockCtrl)

	mockStore.EXPECT().RegisterNode(gomock.Any(), gomock.Any()).Return(errTest).AnyTimes()
	mockStore.EXPECT().GetNodes(gomock.Any(), gomock.Any()).Return(nil, errTest).AnyTimes()

	mockStore.EXPECT().GetNodesByTenantID(gomock.Any(), gomock.Eq(notFoundErrorFilter), gomock.Any()).Return(nil, errNotFound).AnyTimes()
	mockStore.EXPECT().GetNodesByTenantID(gomock.Any(), gomock.Not(gomock.Eq(notFoundErrorFilter)), gomock.Any()).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeName(gomock.Any(), gomock.Any(), gomock.Eq(notFoundErrorFilter)).Return(nil, errNotFound).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeName(gomock.Any(), gomock.Any(), gomock.Not(gomock.Eq(notFoundErrorFilter))).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeID(gomock.Any(), gomock.Any(), gomock.Eq(notFoundErrorFilter)).Return(nil, errNotFound).AnyTimes()
	mockStore.EXPECT().GetNodeByTenantIDAndNodeID(gomock.Any(), gomock.Any(), gomock.Not(gomock.Eq(notFoundErrorFilter))).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().GetNodeByID(gomock.Any(), gomock.Eq("0")).Return(nil, errNotFound).AnyTimes()
	mockStore.EXPECT().GetNodeByID(gomock.Any(), gomock.Not(gomock.Eq("0"))).Return(nil, errTest).AnyTimes()
	mockStore.EXPECT().UpdateNodeByName(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, node *models.Node) error {
			if node.Name == notFoundErrorFilter {
				return errors.NotFoundError("not found error")
			}
			return errTest
		}).AnyTimes()
	mockStore.EXPECT().UpdateNodeByID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, node *models.Node) error {
			if node.ID == "0" {
				return errors.NotFoundError("not found error")
			}
			return errTest
		}).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByName(gomock.Any(), gomock.Eq(notFoundErrorFilter), gomock.Any(), gomock.Any()).Return(errNotFound).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByName(gomock.Any(), gomock.Not(gomock.Eq(notFoundErrorFilter)), gomock.Any(), gomock.Any()).Return(errTest).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByID(gomock.Any(), gomock.Eq("0"), gomock.Any()).Return(errNotFound).AnyTimes()
	mockStore.EXPECT().UpdateBlockPositionByID(gomock.Any(), gomock.Not(gomock.Eq("0")), gomock.Any()).Return(errTest).AnyTimes()
	mockStore.EXPECT().DeleteNodeByName(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, node *models.Node) error {
			if node.Name == notFoundErrorFilter {
				return errors.NotFoundError("not found error")
			}
			return errTest

		}).AnyTimes()
	mockStore.EXPECT().DeleteNodeByID(gomock.Any(), gomock.Eq("0")).Return(errNotFound).AnyTimes()
	mockStore.EXPECT().DeleteNodeByID(gomock.Any(), gomock.Not(gomock.Eq("0"))).Return(errTest).AnyTimes()

	return mockStore
}
