package pg

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

// ContractRegistry is a contract registry based on PostgreSQL
type ContractRegistry struct {
	db *pg.DB
}

// NewContractRegistry creates a new contract registry
func NewContractRegistry(db *pg.DB) *ContractRegistry {
	return &ContractRegistry{db: db}
}

// NewContractRegistryFromPGOptions creates a new pg contract registry
func NewContractRegistryFromPGOptions(opts *pg.Options) *ContractRegistry {
	return NewContractRegistry(pg.Connect(opts))
}

// RegisterContract register a contract including ABI, bytecode and deployed bytecode
func (r *ContractRegistry) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest) (*svc.RegisterContractResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	contract := req.GetContract()

	bytecode, deployedBytecode, abiRaw, err := common.CheckExtractArtifacts(contract)
	if err != nil {
		return nil, err
	}

	name, tagName, err := common.CheckExtractNameTag(contract.Id)
	if err != nil {
		return nil, err
	}

	repositoryModel := &RepositoryModel{
		Name: name,
	}
	_, err = tx.ModelContext(ctx, repositoryModel).
		Column("id").
		Where("name = ?name").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()
	if err != nil {
		log.WithError(err).Debug("could not create repository")
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	d, err := hexutil.Decode(deployedBytecode)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	artifact := &ArtifactModel{
		Abi:              abiRaw,
		Bytecode:         bytecode,
		DeployedBytecode: deployedBytecode,
		Codehash:         hexutil.Encode(crypto.Keccak256(d)),
	}
	_, err = tx.ModelContext(ctx, artifact).
		Column("id").
		Where("abi = ?abi").
		Where("codehash = ?codehash").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()
	if err != nil {
		log.WithError(err).Debug("could not create artifact")
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	tag := &TagModel{
		Name:         tagName,
		RepositoryID: repositoryModel.ID,
		ArtifactID:   artifact.ID,
	}
	_, err = tx.ModelContext(ctx, tag).
		OnConflict("ON CONSTRAINT tags_name_repository_id_key DO UPDATE").
		Set("artifact_id = ?artifact_id").
		Insert()
	if err != nil {
		log.WithError(err).Debug("could not create tag")
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	if abiRaw != "" {
		codeHash := crypto.Keccak256Hash(d)
		contractAbi, err := contract.ToABI()
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}
		methodJSONs, eventJSONs, err := common.ParseJSONABI(abiRaw)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		var methods []*MethodModel
		for _, m := range contractAbi.Methods {
			// Register methods for this bytecode
			method := m
			sel := common.SigHashToSelector(method.ID())
			if deployedBytecode != "" {
				methods = append(methods, &MethodModel{
					Codehash: codeHash.Hex(),
					Selector: sel,
					ABI:      methodJSONs[method.Sig()],
				})
			}
		}
		if methods != nil {
			_, err = tx.ModelContext(ctx, &methods).
				OnConflict("DO NOTHING").
				Insert()
			if err != nil {
				log.WithError(err).Debug("could not insert methods")
				return nil, errors.FromError(err).ExtendComponent(component)
			}
		}

		var events []*EventModel
		for _, e := range contractAbi.Events {
			event := e
			indexedCount := common.GetIndexedCount(event)
			// Register events for this bytecode
			if deployedBytecode != "" {
				events = append(events, &EventModel{
					Codehash:          codeHash.Hex(),
					SigHash:           event.ID().Hex(),
					IndexedInputCount: indexedCount,
					ABI:               eventJSONs[event.Sig()],
				})
			}
		}
		if events != nil {
			_, err = tx.ModelContext(ctx, &events).
				OnConflict("DO NOTHING").
				Insert()
			if err != nil {
				log.WithError(err).Debug("could not insert events")
				return nil, errors.FromError(err).ExtendComponent(component)
			}
		}
	}

	return &svc.RegisterContractResponse{}, tx.Commit()
}

// DeregisterContract remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
func (r *ContractRegistry) DeregisterContract(ctx context.Context, req *svc.DeregisterContractRequest) (*svc.DeregisterContractResponse, error) {
	return nil, nil
}

// DeleteArtifact remove an artifacts based on its BytecodeHash.
func (r *ContractRegistry) DeleteArtifact(ctx context.Context, req *svc.DeleteArtifactRequest) (*svc.DeleteArtifactResponse, error) {
	return nil, nil
}

// GetContract loads a contract
func (r *ContractRegistry) GetContract(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractResponse, error) {
	name, tag, err := common.CheckExtractNameTag(req.GetContractId())
	if err != nil {
		return nil, err
	}

	artifact := &ArtifactModel{}
	err = r.db.ModelContext(ctx, artifact).
		Column("abi", "bytecode", "deployed_bytecode").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS r ON r.id = t.repository_id").
		Where("t.name = ?", tag).
		Where("r.name = ?", name).
		First()
	if err != nil {
		log.WithError(err).Debugf("could not load contract with name: %s and tag: %s", name, tag)
		return nil, errors.StorageError("could not load contract (%v)", err).ExtendComponent(component)
	}

	return &svc.GetContractResponse{
		Contract: &abi.Contract{
			Id: &abi.ContractId{
				Name: name,
				Tag:  tag,
			},
			Abi:              artifact.Abi,
			Bytecode:         artifact.Bytecode,
			DeployedBytecode: artifact.DeployedBytecode,
		},
	}, nil
}

// GetContractABI loads contract ABI
func (r *ContractRegistry) GetContractABI(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractABIResponse, error) {
	name, tag, err := common.CheckExtractNameTag(req.GetContractId())
	if err != nil {
		return nil, err
	}

	artifact := &ArtifactModel{}
	err = r.db.ModelContext(ctx, artifact).
		Column("abi").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS r ON r.id = t.repository_id").
		Where("t.name = ?", tag).
		Where("r.name = ?", name).
		First()
	if err != nil {
		log.WithError(err).Debugf("could not load contract with name: %s and tag: %s", name, tag)
		return nil, errors.StorageError("could not load contract (%v)", err).ExtendComponent(component)
	}

	return &svc.GetContractABIResponse{
		Abi: artifact.Abi,
	}, nil
}

// GetContractBytecode loads contract bytecode
func (r *ContractRegistry) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractBytecodeResponse, error) {
	name, tag, err := common.CheckExtractNameTag(req.GetContractId())
	if err != nil {
		return nil, err
	}

	artifact := &ArtifactModel{}
	err = r.db.ModelContext(ctx, artifact).
		Column("bytecode").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS r ON r.id = t.repository_id").
		Where("t.name = ?", tag).
		Where("r.name = ?", name).
		First()
	if err != nil {
		log.WithError(err).Debugf("could not load contract with name: %s and tag: %s", name, tag)
		return nil, errors.StorageError("could not load contract (%v)", err).ExtendComponent(component)
	}

	return &svc.GetContractBytecodeResponse{
		Bytecode: artifact.Bytecode,
	}, nil
}

// GetContractDeployedBytecode loads contract deployed bytecode
func (r *ContractRegistry) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest) (*svc.GetContractDeployedBytecodeResponse, error) {
	name, tag, err := common.CheckExtractNameTag(req.GetContractId())
	if err != nil {
		return nil, err
	}

	artifact := &ArtifactModel{}
	err = r.db.ModelContext(ctx, artifact).
		Column("deployed_bytecode").
		Join("JOIN tags AS t ON t.artifact_id = artifact_model.id").
		Join("JOIN repositories AS r ON r.id = t.repository_id").
		Where("t.name = ?", tag).
		Where("r.name = ?", name).
		First()
	if err != nil {
		log.WithError(err).Debugf("could not load contract with name: %s and tag: %s", name, tag)
		return nil, errors.StorageError("could not load contract (%v)", err).ExtendComponent(component)
	}

	return &svc.GetContractDeployedBytecodeResponse{
		DeployedBytecode: artifact.DeployedBytecode,
	}, nil
}

// GetMethodsBySelector load method using 4 bytes unique selector and the address of the contract
func (r *ContractRegistry) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest) (*svc.GetMethodsBySelectorResponse, error) {
	method := &MethodModel{}
	err := r.db.ModelContext(ctx, method).
		Column("method_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = method_model.codehash").
		Where("c.chain_id = ?", req.GetAccountInstance().GetChainId()).
		Where("c.address = ?", req.GetAccountInstance().GetAccount()).
		Where("method_model.selector = ?", req.GetSelector()).
		First()
	if err == nil {
		return &svc.GetMethodsBySelectorResponse{
			Method: method.ABI,
		}, nil
	}

	var defaultMethods []*MethodModel
	err = r.db.ModelContext(ctx, &defaultMethods).
		ColumnExpr("DISTINCT abi").
		Where("selector = ?", req.GetSelector()).
		Select()
	if err != nil || len(defaultMethods) == 0 {
		log.WithError(err).Errorf("could not load method: %s", err)
		return nil, errors.NotFoundError("method not found").ExtendComponent(component)
	}

	var methodsABI []string
	for _, m := range defaultMethods {
		methodsABI = append(methodsABI, m.ABI)
	}
	return &svc.GetMethodsBySelectorResponse{
		DefaultMethods: methodsABI,
	}, nil
}

// GetEventsBySigHash load event using event signature hash
func (r *ContractRegistry) GetEventsBySigHash(ctx context.Context, req *svc.GetEventsBySigHashRequest) (*svc.GetEventsBySigHashResponse, error) {
	event := &EventModel{}
	err := r.db.ModelContext(ctx, event).
		Column("event_model.abi").
		Join("JOIN codehashes AS c ON c.codehash = event_model.codehash").
		Where("c.chain_id = ?", req.GetAccountInstance().GetChainId()).
		Where("c.address = ?", req.GetAccountInstance().GetAccount()).
		Where("event_model.sig_hash = ?", req.GetSigHash()).
		Where("event_model.indexed_input_count = ?", req.GetIndexedInputCount()).
		First()
	if err == nil {
		return &svc.GetEventsBySigHashResponse{
			Event: event.ABI,
		}, nil
	}

	var defaultEvents []*EventModel
	err = r.db.ModelContext(ctx, &defaultEvents).
		ColumnExpr("DISTINCT abi").
		Where("sig_hash = ?", req.GetSigHash()).
		Where("indexed_input_count = ?", req.GetIndexedInputCount()).
		Order("abi DESC").
		Select()
	if err != nil || len(defaultEvents) == 0 {
		log.WithError(err).Warn("could not load event")
		return nil, errors.NotFoundError("event not found").ExtendComponent(component)
	}

	var eventsABI []string
	for _, e := range defaultEvents {
		eventsABI = append(eventsABI, e.ABI)
	}
	return &svc.GetEventsBySigHashResponse{
		DefaultEvents: eventsABI,
	}, nil
}

// GetCatalog returns a list of all registered contracts.
func (r *ContractRegistry) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest) (*svc.GetCatalogResponse, error) {
	var names []string
	err := r.db.ModelContext(ctx, (*RepositoryModel)(nil)).
		Column("name").
		OrderExpr("lower(name)").
		Select(&names)
	if err != nil {
		log.WithError(err).Error("could not get Catalog")
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return &svc.GetCatalogResponse{Names: names}, nil
}

// GetTags returns a list of all tags available for a contract name.
func (r *ContractRegistry) GetTags(ctx context.Context, req *svc.GetTagsRequest) (*svc.GetTagsResponse, error) {
	var names []string
	err := r.db.ModelContext(ctx, (*TagModel)(nil)).
		Column("tag_model.name").
		Join("JOIN repositories AS r ON r.id = tag_model.repository_id").
		Where("r.name = ?", req.GetName()).
		OrderExpr("lower(tag_model.name)").
		Select(&names)
	if err != nil {
		log.WithError(err).Error("Could not get Tags")
		return nil, errors.FromError(err).ExtendComponent(component)
	}
	return &svc.GetTagsResponse{Tags: names}, nil
}

// SetAccountCodeHash set the codehash of a contract address for a given chain
func (r *ContractRegistry) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest) (*svc.SetAccountCodeHashResponse, error) {
	codehash := &CodehashModel{
		ChainID:  req.GetAccountInstance().GetChainId(),
		Address:  req.GetAccountInstance().GetAccount(),
		Codehash: req.GetCodeHash(),
	}

	// Execute ORM query
	// If uniqueness constraint is broken then it update the former value
	_, err := r.db.ModelContext(ctx, codehash).
		OnConflict("ON CONSTRAINT codehashes_chain_id_address_key DO UPDATE").
		Set("chain_id = ?chain_id").
		Set("address = ?address").
		Set("codehash = ?codehash").
		Returning("*").
		Insert()
	if err != nil {
		log.WithError(err).Error("Could not store")
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return &svc.SetAccountCodeHashResponse{}, nil
}
