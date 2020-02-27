package client

import (
	"context"
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

var MockChainsSlice = []*types.Chain{
	{
		UUID:                    "test-chain",
		Name:                    "test",
		TenantID:                "test",
		URLs:                    []string{"test"},
		ListenerDepth:           &(&struct{ x uint64 }{0}).x,
		ListenerCurrentBlock:    &(&struct{ x uint64 }{10}).x,
		ListenerStartingBlock:   &(&struct{ x uint64 }{0}).x,
		ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
	},
	{
		UUID:                    "test-chain1",
		Name:                    "test1",
		TenantID:                "test1",
		URLs:                    []string{"test1"},
		ListenerDepth:           &(&struct{ x uint64 }{1}).x,
		ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
		ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
		ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
	},
}

func (s *ClientTestSuite) TestGetChainByUUID() {
	testSuite := []struct {
		server        func() *httptest.Server
		input         string
		testOutput    func(t *testing.T, output *types.Chain, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal(MockChainsSlice[0])
						_, _ = rw.Write(r)
					}),
				)
			},
			MockChainsSlice[0].UUID,
			func(t *testing.T, output *types.Chain, err error) {
				assert.NoError(t, err, "should not get error")
				assert.True(t, reflect.DeepEqual(output, MockChainsSlice[0]), "should be the same chain")
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
			MockChainsSlice[0].UUID,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
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
			MockChainsSlice[0].UUID,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
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
			MockChainsSlice[0].UUID,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			server.Client(),
			&Config{URL: server.URL},
		)

		output, err := s.client.GetChainByUUID(context.Background(), test.input)
		test.testOutput(s.T(), output, err)
		server.Close()
	}
}

func (s *ClientTestSuite) TestGetChainByName() {
	testSuite := []struct {
		server        func() *httptest.Server
		input         string
		testOutput    func(t *testing.T, output *types.Chain, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal(MockChainsSlice)
						_, _ = rw.Write(r)
					}),
				)
			},
			MockChainsSlice[0].Name,
			func(t *testing.T, output *types.Chain, err error) {
				assert.NoError(t, err, "should not get error")
				assert.True(t, reflect.DeepEqual(output, MockChainsSlice[0]), "should be the same chain")
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
			MockChainsSlice[0].Name,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
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
			MockChainsSlice[0].Name,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
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
			MockChainsSlice[0].Name,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal([]*types.Chain{})
						_, _ = rw.Write(r)
					}),
				)
			},
			MockChainsSlice[0].Name,
			func(t *testing.T, output *types.Chain, err error) {
				assert.Error(t, err, "should get error")
				assert.Nil(t, output, "should be nil")
			},
			false,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			server.Client(),
			&Config{URL: server.URL},
		)

		output, err := s.client.GetChainByName(context.Background(), test.input)
		test.testOutput(s.T(), output, err)
		server.Close()
	}
}

func (s *ClientTestSuite) TestGetChains() {
	testSuite := []struct {
		server        func() *httptest.Server
		testOutput    func(t *testing.T, output []*types.Chain, err error)
		expectedError bool
	}{
		{
			func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
						r, _ := json.Marshal(MockChainsSlice)
						_, _ = rw.Write(r)
					}),
				)
			},
			func(t *testing.T, output []*types.Chain, err error) {
				assert.NoError(t, err, "should not get error")
				assert.True(t, reflect.DeepEqual(output, MockChainsSlice), "should be the same chain")
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
			func(t *testing.T, output []*types.Chain, err error) {
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
			func(t *testing.T, output []*types.Chain, err error) {
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
			func(t *testing.T, output []*types.Chain, err error) {
				assert.Error(t, err, "should not get error")
				assert.Nil(t, output, "should be nil")
			},
			true,
		},
	}

	for _, test := range testSuite {
		server := test.server()
		s.client = NewHTTPClient(
			server.Client(),
			&Config{URL: server.URL},
		)

		output, err := s.client.GetChains(context.Background())
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
			server.Client(),
			&Config{URL: server.URL},
		)

		err := s.client.UpdateBlockPosition(context.Background(), "chainUUID", 1)
		test.testOutput(s.T(), err)
		server.Close()
	}
}

func TestClientSuite(t *testing.T) {
	s := new(ClientTestSuite)
	suite.Run(t, s)
}
