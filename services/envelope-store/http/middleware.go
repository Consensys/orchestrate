package http

import (
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	authmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/auth"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
)

type middlewareBuilder middleware.Builder

func newMiddlewareBuilder(jwt, key auth.Checker, multitenancy bool) (middlewareBuilder, error) {
	builder := dynmid.NewBuilder()

	// Add Authentication Middleware
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Auth{}),
		authmid.NewBuilder(jwt, key, multitenancy),
	)

	return builder, nil
}