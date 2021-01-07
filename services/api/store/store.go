package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Store interface {
	Connect(ctx context.Context, conf interface{}) (DB, error)
}

type Agents interface {
	Schedule() ScheduleAgent
	Job() JobAgent
	Log() LogAgent
	Transaction() TransactionAgent
	TransactionRequest() TransactionRequestAgent
	Account() AccountAgent
	Faucet() FaucetAgent
	Artifact() ArtifactAgent
	CodeHash() CodeHashAgent
	Event() EventAgent
	Method() MethodAgent
	Repository() RepositoryAgent
	Tag() TagAgent
}

type DB interface {
	database.DB
	Agents
}

type Tx interface {
	database.Tx
	Agents
}

// Interfaces data agents
type TransactionRequestAgent interface {
	Insert(ctx context.Context, txRequest *models.TransactionRequest) error
	FindOneByIdempotencyKey(ctx context.Context, idempotencyKey string, tenantID string) (*models.TransactionRequest, error)
	FindOneByUUID(ctx context.Context, scheduleUUID string, tenants []string) (*models.TransactionRequest, error)
	Search(ctx context.Context, filters *entities.TransactionRequestFilters, tenants []string) ([]*models.TransactionRequest, error)
}

type ScheduleAgent interface {
	Insert(ctx context.Context, schedule *models.Schedule) error
	FindOneByUUID(ctx context.Context, uuid string, tenants []string) (*models.Schedule, error)
	FindAll(ctx context.Context, tenants []string) ([]*models.Schedule, error)
}

type JobAgent interface {
	Insert(ctx context.Context, job *models.Job) error
	Update(ctx context.Context, job *models.Job) error
	FindOneByUUID(ctx context.Context, uuid string, tenants []string) (*models.Job, error)
	LockOneByUUID(ctx context.Context, uuid string) error
	Search(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*models.Job, error)
}

type LogAgent interface {
	Insert(ctx context.Context, log *models.Log) error
}

type TransactionAgent interface {
	Insert(ctx context.Context, tx *models.Transaction) error
	Update(ctx context.Context, tx *models.Transaction) error
}

type AccountAgent interface {
	Insert(ctx context.Context, account *models.Account) error
	Update(ctx context.Context, account *models.Account) error
	FindOneByAddress(ctx context.Context, address string, tenants []string) (*models.Account, error)
	Search(ctx context.Context, filters *entities.AccountFilters, tenants []string) ([]*models.Account, error)
}

type FaucetAgent interface {
	Insert(ctx context.Context, faucet *models.Faucet) error
	Update(ctx context.Context, faucet *models.Faucet, tenants []string) error
	FindOneByUUID(ctx context.Context, uuid string, tenants []string) (*models.Faucet, error)
	Search(ctx context.Context, filters *entities.FaucetFilters, tenants []string) ([]*models.Faucet, error)
	Delete(ctx context.Context, faucet *models.Faucet, tenants []string) error
}

type ArtifactAgent interface {
	FindOneByABIAndCodeHash(ctx context.Context, abi, codeHash string) (*models.ArtifactModel, error)
	SelectOrInsert(ctx context.Context, artifact *models.ArtifactModel) error
	Insert(ctx context.Context, artifact *models.ArtifactModel) error
	FindOneByNameAndTag(ctx context.Context, name, tag string) (*models.ArtifactModel, error)
}

type CodeHashAgent interface {
	Insert(ctx context.Context, codehash *models.CodehashModel) error
}

type EventAgent interface {
	InsertMultiple(ctx context.Context, events []*models.EventModel) error
	FindOneByAccountAndSigHash(ctx context.Context, chainID, address, sighash string, indexedInputCount uint32) (*models.EventModel, error)
	FindDefaultBySigHash(ctx context.Context, sighash string, indexedInputCount uint32) ([]*models.EventModel, error)
}

type MethodAgent interface {
	InsertMultiple(ctx context.Context, methods []*models.MethodModel) error
	FindOneByAccountAndSelector(ctx context.Context, chainID, address string, selector []byte) (*models.MethodModel, error)
	FindDefaultBySelector(ctx context.Context, selector []byte) ([]*models.MethodModel, error)
}

type RepositoryAgent interface {
	SelectOrInsert(ctx context.Context, repository *models.RepositoryModel) error
	Insert(ctx context.Context, repository *models.RepositoryModel) error
	FindOne(ctx context.Context, name string) (*models.RepositoryModel, error)
	FindAll(ctx context.Context) ([]string, error)
}

type TagAgent interface {
	Insert(ctx context.Context, tag *models.TagModel) error
	FindAllByName(ctx context.Context, name string) ([]string, error)
}
