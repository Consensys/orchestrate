package controllers

import (
	"context"
	"fmt"
	"net/http"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
)

//go:generate swag init --generalInfo builder.go --output ../../../../public/swagger-specs/services/key-manager --parseDependency --parseDepth 2
//go:generate rm ../../../../public/swagger-specs/services/key-manager/docs.go ../../../../public/swagger-specs/services/key-manager/swagger.yaml

// @title Key Management API
// @version 2.0
// @description ConsenSys Codefi Orchestrate Key Management. Enables fine-grained management of cryptographic keys.

// @contact.name Contact ConsenSys Codefi Orchestrate
// @contact.url https://consensys.net/codefi/orchestrate/contact
// @contact.email orchestrate@consensys.net

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

type Builder struct {
	ethereumCtrl *EthereumController
	zksCtrl      *ZKSController
}

func NewBuilder(vault store.Vault, ethUseCases usecases.ETHUseCases, zksUseCases usecases.ZKSUseCases) *Builder {
	return &Builder{
		ethereumCtrl: NewEthereumController(vault, ethUseCases),
		zksCtrl:      NewZKSController(vault, zksUseCases),
	}
}

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}, respModifier func(response *http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.KeyManager)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	b.ethereumCtrl.Append(router)
	b.zksCtrl.Append(router)

	return router, nil
}
