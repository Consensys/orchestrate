package chains

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
