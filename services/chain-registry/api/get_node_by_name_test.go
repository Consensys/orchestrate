package api

import (
	"net/http"
	"strings"
)

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
