package dataagents

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	pg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const contractDAComponent = "data-agents.contract"

type PGContract struct {
	db     pg.DB
	logger *log.Logger
}

type contractQuery struct {
	Name             string `json:"name,omitempty"`
	Tag              string `json:"tag,omitempty"`
	ABI              string `json:"abi,omitempty"`
	Bytecode         string `json:"bytecode,omitempty"`
	DeployedBytecode string `json:"deployed_bytecode,omitempty"`
}

func NewPGContract(db pg.DB) store.ContractAgent {
	return &PGContract{
		db:     db,
		logger: log.NewLogger().SetComponent(contractDAComponent),
	}
}

func (agent *PGContract) FindOneByCodeHash(ctx context.Context, codeHash string) (*entities.Contract, error) {
	qContract := &contractQuery{}
	query := `
SELECT a.abi, a.bytecode, a.deployed_bytecode, r.name as name, t.name as tag 
FROM artifacts a
INNER JOIN tags t ON (a.id = t.artifact_id)
INNER JOIN repositories r ON (r.id = t.repository_id) 
WHERE a.codehash = ?
ORDER BY t.id DESC
LIMIT 1
`
	_, err := agent.db.QueryOneContext(ctx, qContract, query, codeHash)
	if err != nil {
		return nil, pg.ParsePGError(err)
	}

	return parseContract(qContract), nil
}

func (agent *PGContract) FindOneByAddress(ctx context.Context, address string) (*entities.Contract, error) {
	qContract := &contractQuery{}
	query := `
SELECT a.abi, a.bytecode, a.deployed_bytecode, r.name as name, t.name as tag 
FROM artifacts a
INNER JOIN tags t ON (a.id = t.artifact_id)
INNER JOIN repositories r ON (r.id = t.repository_id) 
INNER JOIN codehashes ch ON (ch.codehash = a.codehash) 
WHERE ch.address = ?
LIMIT 1
`
	_, err := agent.db.QueryOneContext(ctx, qContract, query, address)
	if err != nil {
		return nil, pg.ParsePGError(err)
	}

	return parseContract(qContract), nil
}

func parseContract(qContract *contractQuery) *entities.Contract {
	return &entities.Contract{
		Name:             qContract.Name,
		Tag:              qContract.Tag,
		ABI:              qContract.ABI,
		Bytecode:         hexutil.MustDecode(qContract.Bytecode),
		DeployedBytecode: hexutil.MustDecode(qContract.DeployedBytecode),
	}
}
