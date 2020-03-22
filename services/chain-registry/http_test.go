// +build unit

package chainregistry

import (
	"context"
	"testing"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	mockethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

func TestRouterBuilder(t *testing.T) {
	cfg := http.DefaultConfig()
	cfg.API = &traefikstatic.API{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockstore.NewMockChainRegistryStore(ctrl)
	ec := mockethclient.NewMockChainLedgerReader(ctrl)
	jwt := mockauth.NewMockChecker(ctrl)
	key := mockauth.NewMockChecker(ctrl)

	builder, err := NewHTTPBuilder(
		cfg,
		jwt, key,
		true,
		store, ec,
	)
	require.NoError(t, err)

	chains := []*types.Chain{
		{
			UUID:     "0d60a85e-0b90-4482-a14c-108aea2557aa",
			Name:     "testChain",
			TenantID: "testTenantId",
			URLs: []string{
				"http://testURL1.com",
				"http://testURL2.com",
			},
		},
	}

	dynCfgs := map[string]interface{}{
		InternalProvider:    NewInternalConfig(cfg),
		ChainsProxyProvider: NewChainsProxyConfig(chains),
	}

	dyncCfg := dynamic.Merge(dynCfgs)
	_, err = builder.Build(
		context.Background(),
		[]string{http.DefaultHTTPEntryPoint, http.DefaultMetricsEntryPoint},
		dyncCfg,
	)
	require.NoError(t, err, "Build router should not error")
}
