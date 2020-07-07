// +build unit

package chains

import (
	"encoding/json"
	"net/http"
)

const (
	name                     = "testName"
	defaultResult            = "{\"uuid\":\"uuid\",\"name\":\"testName\",\"tenantID\":\"_\",\"urls\":[\"http://test.com\"],\"chainID\":\"888\",\"listenerDepth\":0,\"listenerCurrentBlock\":\"666\",\"listenerStartingBlock\":\"666\",\"listenerBackOffDuration\":\"1s\",\"listenerExternalTxEnabled\":false,\"createdAt\":null}\n"
	defaultMultitenantResult = "{\"uuid\":\"uuid\",\"name\":\"testName\",\"tenantID\":\"tenantID\",\"urls\":[\"http://test.com\"],\"chainID\":\"888\",\"listenerDepth\":0,\"listenerCurrentBlock\":\"666\",\"listenerStartingBlock\":\"666\",\"listenerBackOffDuration\":\"1s\",\"listenerExternalTxEnabled\":false,\"createdAt\":null}\n"
)

var urls = []string{"http://test.com"}

var postChainTests = []HTTPRouteTests{
	{
		name:       "TestPostChain200",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			listenerDepth := uint64(1)
			listenerFromBlock := "666"
			listenerBackOffDuration := "1s"
			listenerExternalTxEnabled := true

			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
				Listener: &ListenerPostRequest{
					Depth:             &listenerDepth,
					FromBlock:         &listenerFromBlock,
					BackOffDuration:   &listenerBackOffDuration,
					ExternalTxEnabled: &listenerExternalTxEnabled,
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return "{\"uuid\":\"uuid\",\"name\":\"testName\",\"tenantID\":\"_\",\"urls\":[\"http://test.com\"],\"chainID\":\"888\",\"listenerDepth\":1,\"listenerCurrentBlock\":\"666\",\"listenerStartingBlock\":\"666\",\"listenerBackOffDuration\":\"1s\",\"listenerExternalTxEnabled\":true,\"createdAt\":null}\n"
		},
	},
	{
		name:       "TestPostTesseraChain200",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			listenerDepth := uint64(1)
			listenerFromBlock := "500"
			listenerBackOffDuration := "1s"
			listenerExternalTxEnabled := true

			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
				Listener: &ListenerPostRequest{
					Depth:             &listenerDepth,
					FromBlock:         &listenerFromBlock,
					BackOffDuration:   &listenerBackOffDuration,
					ExternalTxEnabled: &listenerExternalTxEnabled,
				},
				PrivateTxManager: &PrivateTxManagerRequest{
					URL:  "http://tessera.com",
					Type: "Tessera",
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return "{\"uuid\":\"uuid\",\"name\":\"testName\",\"tenantID\":\"_\",\"urls\":[\"http://test.com\"],\"chainID\":\"888\",\"listenerDepth\":1,\"listenerCurrentBlock\":\"500\",\"listenerStartingBlock\":\"500\",\"listenerBackOffDuration\":\"1s\",\"listenerExternalTxEnabled\":true,\"createdAt\":null,\"privateTxManagers\":[{\"UUID\":\"uuid\",\"ChainUUID\":\"uuid\",\"url\":\"http://tessera.com\",\"type\":\"Tessera\",\"CreatedAt\":null}]}\n"
		},
	},
	{
		name:       "TestPostChain200 Listener is Nil",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return defaultResult
		},
	},
	{
		name:       "TestPostChain200 FromBlock is Nil",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name:     name,
				URLs:     urls,
				Listener: &ListenerPostRequest{},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return defaultResult
		},
	},
	{
		name:       "TestPostChain200 FromBlock is latest",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			listenerFromBlock := latestBlockStr

			body, _ := json.Marshal(&PostRequest{
				Name: name,
				URLs: urls,
				Listener: &ListenerPostRequest{
					FromBlock: &listenerFromBlock,
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return defaultResult
		},
	},

	{
		name:       "TestPostChain400WithTwiceSameURL",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: name,
				URLs: []string{"http://test.com", "http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			return "{\"message\":\"42400@encoding.json: invalid body, with: field validation for 'URLs' failed on the 'unique' tag\"}\n"
		},
	},
	{
		name:       "TestPostChain400WrongURL",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: name,
				URLs: []string{"test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			return expectedNotUniqueURLsError
		},
	},
	{
		name:                "TestPostChain400WrongBody",
		chainAgent:          UseMockChainRegistry,
		httpMethod:          http.MethodPost,
		path:                "/chains",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedUnknownBodyError },
	},
	{
		name:       "TestPostChain500",
		chainAgent: UseErrorChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {

			body, _ := json.Marshal(&PostRequest{
				Name: name,
				URLs: urls,
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return "{\"message\":\"FF000@use-cases.register-chain: test error\"}\n" },
	},
	{
		name:                "TestPostChain400WrongBodyTessera",
		chainAgent:          UseMockChainRegistry,
		httpMethod:          http.MethodPost,
		path:                "/chains",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedUnknownBodyError },
	},
	{
		name:       "TestPostTesseraChain500",
		chainAgent: UseErrorChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			listenerDepth := uint64(1)
			listenerFromBlock := "500"
			listenerBackOffDuration := "1s"
			listenerExternalTxEnabled := true

			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
				Listener: &ListenerPostRequest{
					Depth:             &listenerDepth,
					FromBlock:         &listenerFromBlock,
					BackOffDuration:   &listenerBackOffDuration,
					ExternalTxEnabled: &listenerExternalTxEnabled,
				},
				PrivateTxManager: &PrivateTxManagerRequest{
					URL:  "http://tessera.com",
					Type: "InvalidType",
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedErrorInvalidManagerType },
	},
}

var postChainTestsMultitenant = []HTTPRouteTests{
	{
		name:       "TestPostChain200Multitenant",
		chainAgent: UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: name,
				URLs: urls,
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody: func() string {
			return defaultMultitenantResult
		},
	},
}
