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
