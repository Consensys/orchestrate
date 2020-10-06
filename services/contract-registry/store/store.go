package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

//go:generate mockgen -source=store.go -destination=mock/mock.go -package=mock

// Interfaces data agents
type ContractDataAgent interface {
	Insert(
		ctx context.Context,
		name, tagName, abiRaw, bytecode, deployedBytecode, codeHash string,
		methods *[]*models.MethodModel,
		events *[]*models.EventModel,
	) error
}

type ArtifactDataAgent interface {
	SelectOrInsert(ctx context.Context, artifact *models.ArtifactModel) error
	FindOneByNameAndTag(ctx context.Context, name, tag string) (*models.ArtifactModel, error)
}

type CodeHashDataAgent interface {
	Insert(ctx context.Context, codehash *models.CodehashModel) error
}

type EventDataAgent interface {
	InsertMultiple(ctx context.Context, events *[]*models.EventModel) error
	FindOneByAccountAndSigHash(ctx context.Context, account *common.AccountInstance, sighash string, indexedInputCount uint32) (*models.EventModel, error)
	FindDefaultBySigHash(ctx context.Context, sighash string, indexedInputCount uint32) ([]*models.EventModel, error)
}

type MethodDataAgent interface {
	InsertMultiple(ctx context.Context, methods *[]*models.MethodModel) error
	FindOneByAccountAndSelector(ctx context.Context, account *common.AccountInstance, selector []byte) (*models.MethodModel, error)
	FindDefaultBySelector(ctx context.Context, selector []byte) ([]*models.MethodModel, error)
}

type RepositoryDataAgent interface {
	SelectOrInsert(ctx context.Context, repository *models.RepositoryModel) error
	FindAll(ctx context.Context) ([]string, error)
}

type TagDataAgent interface {
	Insert(ctx context.Context, tag *models.TagModel) error
	FindAllByName(ctx context.Context, name string) ([]string, error)
}
