package client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type ClientTestSuite struct {
	suite.Suite
	client *HTTPClient
}

var MockNodesSlice = []*types.Node{
	{
		ID:                      "test-node",
		Name:                    "test",
		TenantID:                "test",
		URLs:                    []string{"test"},
		ListenerDepth:           0,
		ListenerBlockPosition:   0,
		ListenerFromBlock:       0,
		ListenerBackOffDuration: "0s",
	},
	{
		ID:                      "test-node1",
		Name:                    "test1",
		TenantID:                "test1",
		URLs:                    []string{"test1"},
		ListenerDepth:           1,
		ListenerBlockPosition:   1,
		ListenerFromBlock:       1,
		ListenerBackOffDuration: "1s",
	},
}

func (s *ClientTestSuite) TestGetNodeByID() {
	testSuite := []struct {
		server        func() *httptest.Server
		input         string
		testOutput    func(t *testing.T, output *types.Node, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal(MockNodesSlice[0])
						_, _ = rw.Write(r)
					}),
				)
			},
			MockNodesSlice[0].ID,
			func(t *testing.T, output *types.Node, err error) {
				assert.NoError(t, err, "should not get error")
				assert.True(t, reflect.DeepEqual(output, MockNodesSlice[0]), "should be the same node")
			},
			false,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						http.Error(rw, "error", 500)
					}),
				)
			},
			MockNodesSlice[0].ID,
			func(t *testing.T, output *types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						_, _ = rw.Write([]byte(`test`))
					}),
				)
			},
			MockNodesSlice[0].ID,
			func(t *testing.T, output *types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
		{
			func() *httptest.Server {
				h := httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}),
				)
				h.Close()
				return h
			},
			MockNodesSlice[0].ID,
			func(t *testing.T, output *types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			*server.Client(),
			Config{URL: server.URL},
		)

		output, err := s.client.GetNodeByID(test.input)
		test.testOutput(s.T(), output, err)
		server.Close()
	}
}

func (s *ClientTestSuite) TestGetNodes() {
	testSuite := []struct {
		server        func() *httptest.Server
		testOutput    func(t *testing.T, output []*types.Node, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal(MockNodesSlice)
						_, _ = rw.Write(r)
					}),
				)
			},
			func(t *testing.T, output []*types.Node, err error) {
				assert.NoError(t, err, "should not get error")
				assert.True(t, reflect.DeepEqual(output, MockNodesSlice), "should be the same node")
			},
			false,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						http.Error(rw, "error", 500)
					}),
				)
			},
			func(t *testing.T, output []*types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						_, _ = rw.Write([]byte(`test`))
					}),
				)
			},
			func(t *testing.T, output []*types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
		{
			func() *httptest.Server {
				h := httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}),
				)
				h.Close()
				return h
			},
			func(t *testing.T, output []*types.Node, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			*server.Client(),
			Config{URL: server.URL},
		)

		output, err := s.client.GetNodes()
		test.testOutput(s.T(), output, err)
		server.Close()
	}
}

func (s *ClientTestSuite) TestUpdateBlockPosition() {
	testSuite := []struct {
		server        func() *httptest.Server
		testOutput    func(t *testing.T, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						_, _ = rw.Write([]byte(`{}`))
					}),
				)
			},
			func(t *testing.T, err error) {
				assert.NoError(t, err, "should not get error")
			},
			false,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						http.Error(rw, "error", 500)
					}),
				)
			},
			func(t *testing.T, err error) {
				assert.Error(t, err, "should not get error")
			},
			true,
		},
		{
			func() *httptest.Server {
				h := httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}),
				)
				h.Close()
				return h
			},
			func(t *testing.T, err error) {
				assert.Error(t, err, "should not get error")
			},
			true,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			*server.Client(),
			Config{URL: server.URL},
		)

		err := s.client.UpdateBlockPosition("nodeID", 1)
		test.testOutput(s.T(), err)
		server.Close()
	}
}

func TestClientSuite(t *testing.T) {
	s := new(ClientTestSuite)
	suite.Run(t, s)
}
