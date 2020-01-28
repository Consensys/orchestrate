package memory

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type ChainRegistry struct {
	ChainsByUUID  map[string]*types.Chain
	ChainsByNames map[string]map[string]*types.Chain
}

// NewChainRegistry creates a new chain registry
func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{
		ChainsByUUID:  make(map[string]*types.Chain),
		ChainsByNames: make(map[string]map[string]*types.Chain),
	}
}

func (r *ChainRegistry) RegisterChain(_ context.Context, chain *types.Chain) error {
	chain.SetDefault()

	if !chain.IsValid() {
		return errors.FromError(fmt.Errorf("invalid chain")).ExtendComponent(component)
	}

	if r.ChainsByNames[chain.TenantID] == nil {
		r.ChainsByNames[chain.TenantID] = make(map[string]*types.Chain)
	}

	if r.ChainsByNames[chain.TenantID][chain.Name] != nil {
		return errors.FromError(fmt.Errorf("chain tenantID=%s name=%s already exitst", chain.TenantID, chain.Name)).ExtendComponent(component)
	}

	r.ChainsByNames[chain.TenantID][chain.Name] = chain
	r.ChainsByUUID[chain.UUID] = chain
	return nil
}

func (r *ChainRegistry) GetChains(_ context.Context, filters map[string]string) ([]*types.Chain, error) {
	// TODO: implement filters

	chains := make([]*types.Chain, 0)

	for _, chain := range r.ChainsByUUID {
		chains = append(chains, chain)
	}

	return chains, nil
}

func (r *ChainRegistry) GetChainsByTenantID(_ context.Context, tenantID string, filters map[string]string) ([]*types.Chain, error) {
	// TODO: implement filters

	chains := make([]*types.Chain, 0)

	if tenantChains, ok := r.ChainsByNames[tenantID]; ok {
		for _, chain := range tenantChains {

			chains = append(chains, chain)
		}
	} else {
		return nil, errors.NotFoundError("unknown tenantID=%s", tenantID).ExtendComponent(component)
	}

	return chains, nil
}

func (r *ChainRegistry) GetChainByTenantIDAndName(_ context.Context, tenantID, name string) (*types.Chain, error) {
	if _, ok := r.ChainsByNames[tenantID]; ok {
		if chain, ok := r.ChainsByNames[tenantID][name]; ok {
			return chain, nil
		}
		return nil, errors.NotFoundError("unknown chain with tenantID=%s and name=%s", name, tenantID).ExtendComponent(component)
	}

	return nil, errors.NotFoundError("unknown chain with tenantID=%s", tenantID).ExtendComponent(component)
}

func (r *ChainRegistry) GetChainByTenantIDAndUUID(ctx context.Context, tenantID, id string) (*types.Chain, error) {
	if _, ok := r.ChainsByUUID[id]; !ok || r.ChainsByUUID[id].TenantID != tenantID {
		return nil, errors.FromError(fmt.Errorf("unknown chain UUID=%s", id)).ExtendComponent(component)
	}

	return r.ChainsByUUID[id], nil
}

func (r *ChainRegistry) GetChainByUUID(_ context.Context, id string) (*types.Chain, error) {
	if _, ok := r.ChainsByUUID[id]; !ok {
		return nil, errors.FromError(fmt.Errorf("unknown chain UUID=%s", id)).ExtendComponent(component)
	}

	return r.ChainsByUUID[id], nil
}

func (r *ChainRegistry) UpdateChainByName(_ context.Context, chain *types.Chain) error {
	chainToUpdate, ok := r.ChainsByNames[chain.TenantID][chain.Name]
	if !ok {
		return errors.NotFoundError("no chain found with tenantID %s and name %s", chain.TenantID, chain.Name).ExtendComponent(component)
	}

	chainToUpdate.Name = chain.Name
	chainToUpdate.URLs = chain.URLs
	chainToUpdate.ListenerBackOffDuration = chain.ListenerBackOffDuration
	chainToUpdate.ListenerBlockPosition = chain.ListenerBlockPosition
	chainToUpdate.ListenerFromBlock = chain.ListenerFromBlock
	chainToUpdate.ListenerDepth = chain.ListenerDepth
	currentTime := time.Now()
	chainToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByName(_ context.Context, name, tenantID string, blockPosition int64) error {
	chainToUpdate, ok := r.ChainsByNames[tenantID][name]
	if !ok {
		return errors.NotFoundError("no chain found with tenantID %s and name %s", tenantID, name).ExtendComponent(component)
	}

	*chainToUpdate.ListenerBlockPosition = blockPosition
	currentTime := time.Now()
	chainToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateChainByUUID(_ context.Context, chain *types.Chain) error {
	chainToUpdate, ok := r.ChainsByUUID[chain.UUID]
	if !ok {
		return errors.NotFoundError("no chain found with id %s", chain.UUID).ExtendComponent(component)
	}

	chainToUpdate.Name = chain.Name
	chainToUpdate.URLs = chain.URLs
	chainToUpdate.ListenerBackOffDuration = chain.ListenerBackOffDuration
	chainToUpdate.ListenerBlockPosition = chain.ListenerBlockPosition
	chainToUpdate.ListenerFromBlock = chain.ListenerFromBlock
	chainToUpdate.ListenerDepth = chain.ListenerDepth
	currentTime := time.Now()
	chainToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByUUID(_ context.Context, id string, blockPosition int64) error {
	chainToUpdate, ok := r.ChainsByUUID[id]
	if !ok {
		return errors.NotFoundError("no chain found with id %s", id).ExtendComponent(component)
	}

	*chainToUpdate.ListenerBlockPosition = blockPosition
	currentTime := time.Now()
	chainToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) DeleteChainByName(_ context.Context, chain *types.Chain) error {
	if _, ok := r.ChainsByNames[chain.TenantID]; !ok {
		return errors.NotFoundError("no chain found with tenant_id=%s", chain.TenantID).ExtendComponent(component)
	}

	if _, ok := r.ChainsByNames[chain.TenantID][chain.Name]; !ok {
		return errors.NotFoundError("no chain found with tenant_id=%s and name=%s", chain.TenantID, chain.Name).ExtendComponent(component)
	}

	delete(r.ChainsByUUID, r.ChainsByNames[chain.TenantID][chain.Name].UUID)
	delete(r.ChainsByNames[chain.TenantID], chain.Name)
	return nil
}

func (r *ChainRegistry) DeleteChainByUUID(_ context.Context, id string) error {
	if _, ok := r.ChainsByUUID[id]; !ok {
		return errors.NotFoundError("no chain found with id=%s", id).ExtendComponent(component)
	}

	delete(r.ChainsByNames[r.ChainsByUUID[id].TenantID], r.ChainsByUUID[id].Name)
	delete(r.ChainsByUUID, id)

	return nil
}
