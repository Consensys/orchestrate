package chainregistry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

//go:generate swag init --dir . --generalInfo chain_registry.go --output ../../../public/swagger-specs/types/chain-registry
//go:generate rm ../../../public/swagger-specs/types/chain-registry/docs.go ../../../public/swagger-specs/types/chain-registry/swagger.yaml

// @title Chain Registry API
// @version 2.0
// @description PegaSys Orchestrate Chain Registry API. Enables dynamic management of chains

// @contact.name Contact PegaSys Orchestrate
// @contact.url https://pegasys.tech/contact/
// @contact.email support@pegasys.tech

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization

type Builder struct {
	store store.ChainRegistryStore
	ec    ethclient.ChainLedgerReader
}

func NewBuilder(s store.ChainRegistryStore, ec ethclient.ChainLedgerReader) *Builder {
	return &Builder{
		store: s,
		ec:    ec,
	}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Chains)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	chains.New(b.store, b.ec).Append(router)
	faucets.New(b.store).Append(router)

	return router, nil
}
