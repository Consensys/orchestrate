package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

const (
	expectedInternalServerErrorBody  = "{\"message\":\"test error\"}\n"
	expectedNotFoundErrorBody        = "{\"message\":\"DB200@: not found error\"}\n"
	expectedInvalidIDErrorBody       = "{\"message\":\"invalid ID format\"}\n"
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
	store               models.ChainRegistryStore
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
		postNodeTests,
	}

	for _, tests := range testsSuite {
		for _, tt := range tests {
			tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel() // marks each test case as capable of running in parallel with each other
				c := tt.store
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

type RouteTestSuite struct {
	suite.Suite
	Router *mux.Router
}

func (s *RouteTestSuite) UseMockChainRegistry() {
	c := &MockChainRegistry{}
	r := mux.NewRouter()
	NewHandler(c).Append(r)
	s.Router = r
}

func (s *RouteTestSuite) UseErrorChainRegistry() {
	c := &ErrorChainRegistry{}
	r := mux.NewRouter()
	NewHandler(c).Append(r)
	s.Router = r
}

type MockChainRegistry struct{}

func (e *MockChainRegistry) RegisterNode(_ context.Context, node *models.Node) error {
	node.ID = 1
	node.Name = "nodeName1"
	node.TenantID = "tenantID1"
	node.URLs = []string{"testUrl1", "testUrl2"}
	node.ListenerDepth = 1
	node.ListenerBlockPosition = 1
	node.ListenerFromBlock = 1
	node.ListenerBackOffDuration = "1s"
	return nil
}

func (e *MockChainRegistry) GetNodes(_ context.Context) ([]*models.Node, error) {
	return []*models.Node{}, nil
}
func (e *MockChainRegistry) GetNodesByTenantID(_ context.Context, _ string) ([]*models.Node, error) {
	return []*models.Node{}, nil
}
func (e *MockChainRegistry) GetNodeByName(_ context.Context, _, _ string) (*models.Node, error) {
	return &models.Node{}, nil
}
func (e *MockChainRegistry) GetNodeByID(_ context.Context, _ int) (*models.Node, error) {
	return &models.Node{}, nil
}
func (e *MockChainRegistry) UpdateNodeByName(_ context.Context, _ *models.Node) error { return nil }
func (e *MockChainRegistry) UpdateNodeByID(_ context.Context, _ *models.Node) error   { return nil }
func (e *MockChainRegistry) DeleteNodeByName(_ context.Context, _ *models.Node) error { return nil }
func (e *MockChainRegistry) DeleteNodeByID(_ context.Context, _ int) error            { return nil }

type ErrorChainRegistry struct{}

var errTest = fmt.Errorf("test error")

func (e *ErrorChainRegistry) RegisterNode(_ context.Context, node *models.Node) error {
	return errTest
}
func (e *ErrorChainRegistry) GetNodes(_ context.Context) ([]*models.Node, error) {
	return nil, errTest
}
func (e *ErrorChainRegistry) GetNodesByTenantID(_ context.Context, tenantID string) ([]*models.Node, error) {
	if tenantID == notFoundErrorFilter {
		return nil, errors.NotFoundError("not found error")
	}
	return nil, errTest
}
func (e *ErrorChainRegistry) GetNodeByName(_ context.Context, _, name string) (*models.Node, error) {
	if name == notFoundErrorFilter {
		return nil, errors.NotFoundError("not found error")
	}
	return nil, errTest
}
func (e *ErrorChainRegistry) GetNodeByID(_ context.Context, id int) (*models.Node, error) {
	if id == 0 {
		return nil, errors.NotFoundError("not found error")
	}
	return nil, errTest
}
func (e *ErrorChainRegistry) UpdateNodeByName(_ context.Context, node *models.Node) error {
	if node.Name == notFoundErrorFilter {
		return errors.NotFoundError("not found error")
	}
	return errTest
}
func (e *ErrorChainRegistry) UpdateNodeByID(_ context.Context, node *models.Node) error {
	if node.ID == 0 {
		return errors.NotFoundError("not found error")
	}
	return errTest
}
func (e *ErrorChainRegistry) DeleteNodeByName(_ context.Context, node *models.Node) error {
	if node.Name == notFoundErrorFilter {
		return errors.NotFoundError("not found error")
	}
	return errTest
}
func (e *ErrorChainRegistry) DeleteNodeByID(_ context.Context, id int) error {
	if id == 0 {
		return errors.NotFoundError("not found error")
	}
	return errTest
}

func TestRouteSuite(t *testing.T) {
	s := new(RouteTestSuite)
	suite.Run(t, s)
}
