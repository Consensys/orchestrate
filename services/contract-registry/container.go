package contractregistry

import (
	"github.com/go-pg/pg/v9"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/controllers"
	dataagents "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/data-agents"
)

func initializeController(db *pg.DB) *controllers.ContractRegistryController {
	// Data agents
	repositoryDA := dataagents.NewPGRepository(db)
	artifactDA := dataagents.NewPGArtifact(db)
	codeHashDA := dataagents.NewPGCodeHash(db)
	eventDA := dataagents.NewPGEvent(db)
	methodDA := dataagents.NewPGMethod(db)
	tagDA := dataagents.NewPGTag(db)
	contractDA := dataagents.NewPGContract(db, repositoryDA, artifactDA, tagDA, methodDA, eventDA)

	// Use cases
	registerContractUC := usecases.NewRegisterContract(contractDA)
	getContractUC := usecases.NewGetContract(artifactDA)
	getMethodsUC := usecases.NewGetMethods(methodDA)
	getEventsUC := usecases.NewGetEvents(eventDA)
	getCatalogUC := usecases.NewGetCatalog(repositoryDA)
	getTagsUC := usecases.NewGetTags(tagDA)
	setCodehashUC := usecases.NewSetCodeHash(codeHashDA)

	return controllers.NewContractRegistryController(
		registerContractUC,
		getContractUC,
		getMethodsUC,
		getEventsUC,
		getCatalogUC,
		getTagsUC,
		setCodehashUC)
}
