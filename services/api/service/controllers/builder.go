package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSys/orchestrate/services/key-manager/client"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/gorilla/mux"
)

//go:generate swag init --generalInfo builder.go --parseDepth 100 --output ../../../../public/swagger-specs/services/api --parseDependency --parseDepth 3
//go:generate rm ../../../../public/swagger-specs/services/api/docs.go ../../../../public/swagger-specs/services/api/swagger.yaml

// @title Orchestrate API
// @version 2.0
// @description ConsenSys Codefi Orchestrate API. Enables dynamic management of transactions, identities, chains, faucets and contracts.
// @description Transaction Requests are an abstraction over schedules and jobs representing one or more transactions executed on the Blockchain network
// @description Schedules are ordered lists of jobs executed in a predefined sequence
// @description Jobs represent single transaction flows executed on the Blockchain network
// @description Chains represent list of endpoints pointing to a Blockchain network
// @description Faucets represent funded accounts (holding ETH) linked to specific chains, allowed to fund newly created accounts automatically for them to be able to send transactions.
// @description Accounts represent Ethereum accounts (private keys). By usage of the generated cryptographic key pair, accounts can be used to sign/verify and to encrypt/decrypt messages.
// @description Contracts represent Solidity contracts management.

// @contact.name Contact ConsenSys Codefi Orchestrate
// @contact.url https://consensys.net/codefi/orchestrate/contact
// @contact.email orchestrate@consensys.net

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization

type Builder struct {
	txCtrl        *TransactionsController
	schedulesCtrl *SchedulesController
	jobsCtrl      *JobsController
	accountsCtrl  *AccountsController
	faucetsCtrl   *FaucetsController
	chainsCtrl    *ChainsController
	contractsCtrl *ContractsController
}

func NewBuilder(ucs usecases.UseCases, keyManagerClient client.KeyManagerClient) *Builder {
	return &Builder{
		txCtrl:        NewTransactionsController(ucs),
		schedulesCtrl: NewSchedulesController(ucs),
		jobsCtrl:      NewJobsController(ucs),
		accountsCtrl:  NewAccountsController(ucs, keyManagerClient),
		faucetsCtrl:   NewFaucetsController(ucs),
		chainsCtrl:    NewChainsController(ucs),
		contractsCtrl: NewContractsController(ucs),
	}
}

func (b *Builder) Build(_ context.Context, _ string, configuration interface{}, _ func(response *http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.API)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	b.txCtrl.Append(router)
	b.schedulesCtrl.Append(router)
	b.jobsCtrl.Append(router)
	b.accountsCtrl.Append(router)
	b.faucetsCtrl.Append(router)
	b.chainsCtrl.Append(router)
	b.contractsCtrl.Append(router)

	return router, nil
}
