package api

import (
	"net/http"
	"strings"
)

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
