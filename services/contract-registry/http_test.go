// +build unit

package contractregistry

import (
	"context"
	"testing"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	mockregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/mock"
)

func TestRouterBuilder(t *testing.T) {
	cfg := http.DefaultConfig()
	cfg.API = &traefikstatic.API{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwt := mockauth.NewMockChecker(ctrl)
	key := mockauth.NewMockChecker(ctrl)
	s := mockregistry.NewMockContractRegistryServer(ctrl)

	builder, err := NewHTTPBuilder(
		cfg,
		jwt, key,
		true,
		s,
	)
	require.NoError(t, err)

	dynCfgs := map[string]interface{}{
		InternalProvider: NewInternalConfig(cfg),
	}
	dynCfg := dynamic.Merge(dynCfgs)
	_, err = builder.Build(
		context.Background(),
		[]string{http.DefaultHTTPEntryPoint, http.DefaultMetricsEntryPoint},
		dynCfg,
	)
	require.NoError(t, err, "Build router should not error")
}
