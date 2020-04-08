// +build unit

package http

import (
	"context"
	"testing"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/mock"
	watcher "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/configwatcher"
)

func TestRouterBuilder(t *testing.T) {
	cfg := http.DefaultConfig()
	cfg.API = &traefikstatic.API{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwt := mockauth.NewMockChecker(ctrl)
	key := mockauth.NewMockChecker(ctrl)
	s := mockstore.NewMockEnvelopeStoreServer(ctrl)

	builder, err := NewRouterBuilder(
		s,
		cfg,
		jwt, key,
		true,
	)
	require.NoError(t, err)

	dyncfg := watcher.NewConfig(cfg, &configwatcher.Config{})
	dynCfgs := map[string]interface{}{
		watcher.InternalProviderName: dyncfg.DynamicCfg(),
	}
	dynCfg := dynamic.Merge(dynCfgs)
	_, err = builder.Build(
		context.Background(),
		[]string{http.DefaultHTTPEntryPoint, http.DefaultMetricsEntryPoint},
		dynCfg,
	)
	require.NoError(t, err, "Build router should not error")
}
