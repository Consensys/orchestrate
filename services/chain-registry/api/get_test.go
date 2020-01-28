package api

import (
	"net/http"
)

var getChainsByUUIDTests = []HTTPRouteTests{
	{
		name:                "TestGetChainByUUID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestGetChainByUUID404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains/0",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetChainByUUID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getChainsTests = []HTTPRouteTests{
	{
		name:                "TestGetChains200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetChains200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains?name=test",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetChainByUUID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/chains",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getChainsByTenantIDTests = []HTTPRouteTests{
	{
		name:                "TestGetChainsByTenantID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetChainsByTenantID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains?uuid=test",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetChainsByTenantID404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/notFoundError/chains",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetChainsByTenantID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getChainByNameTests = []HTTPRouteTests{
	{
		name:                "TestGetChainByName200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains/testChainName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestGetChainByName404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains/notFoundError",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetChainByName500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodGet,
		path:                "/testTenantID/chains/testChainName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
