package api

import (
	"net/http"
	"strings"
)

var deleteNodeByIDTests = []HTTPRouteTests{
	{
		name:                "TestDeleteNodeByID200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(deleteNodeByIDPath, "{"+nodeIDPath+"}", "1"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteNodeByID404",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(deleteNodeByIDPath, "{"+nodeIDPath+"}", "0"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteNodeByID500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(deleteNodeByIDPath, "{"+nodeIDPath+"}", "1"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var deleteNodesByNameTests = []HTTPRouteTests{
	{
		name:                "TestDeleteNodeByName200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "testNodeName"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteNodeByName404",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "notFoundError"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteNodeByName500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodDelete,
		path:                strings.ReplaceAll(strings.ReplaceAll(getNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "testNodeName"),
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
