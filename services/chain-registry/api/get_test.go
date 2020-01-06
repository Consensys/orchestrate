package api

import (
	"net/http"
	"strings"
)

var getNodesByIDTests = []HTTPRouteTests{
	{
		name:                "TestGetNodeByID200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodeByIDPath, "{"+nodeIDPath+"}", "1"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestGetNodeByID404",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodeByIDPath, "{"+nodeIDPath+"}", "0"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetNodeByID500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodeByIDPath, "{"+nodeIDPath+"}", "1"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getNodesTests = []HTTPRouteTests{
	{
		name:                "TestGetNodes200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                getNodesPath,
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetNodeByID500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                getNodesPath,
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getNodesByTenantIDTests = []HTTPRouteTests{
	{
		name:                "TestGetNodesByTenantID200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodesByTenantIDPath, "{"+tenantIDPath+"}", "testTenantID"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetNodesByTenantID404",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodesByTenantIDPath, "{"+tenantIDPath+"}", "notFoundError"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetNodesByTenantID500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(getNodesByTenantIDPath, "{"+tenantIDPath+"}", "testTenantID"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var getNodesByNameTests = []HTTPRouteTests{
	{
		name:                "TestGetNodeByName200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{tenantID}", "testTenantID"), "{nodeName}", "testNodeName"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestGetNodeByName404",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{tenantID}", "testTenantID"), "{nodeName}", "notFoundError"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestGetNodeByName500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{tenantID}", "testTenantID"), "{nodeName}", "testNodeName"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
