package controllers

import (
	"context"
	"fmt"
	"net/http"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

//go:generate swag init --dir . --generalInfo builder.go --output ../../../../public/swagger-specs/types/transaction-scheduler
//go:generate rm ../../../../public/swagger-specs/types/transaction-scheduler/docs.go ../../../../public/swagger-specs/types/transaction-scheduler/swagger.yaml

// @title Transaction Scheduler API
// @version 2.0
// @description PegaSys Orchestrate Transaction API. Enables dynamic management of transactions

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
	usecases *usecases.UseCases
}

func NewBuilder(uc *usecases.UseCases) *Builder {
	return &Builder{
		usecases: uc,
	}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(response *http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Transactions)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	NewJobsController(b.usecases).Append(router)
	// TODO: Add schedules and transactions

	return router, nil
}
