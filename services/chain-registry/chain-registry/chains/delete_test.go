// +build unit

package chains

import (
	"net/http"
)

var deleteChainByUUIDTests = []HTTPRouteTests{
	{
		name:                "TestDeleteChainByUUID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNoContent,
		expectedContentType: "",
		expectedBody:        func() string { return "" },
	},
	{
		name:                "TestDeleteChainByUUID404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/0",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteChainByUUID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
