package memory

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type ChainRegistry struct {
	NodesByID    map[string]*types.Node
	NodesByNames map[string]map[string]*types.Node
}

// NewChainRegistry creates a new chain registry
func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{
		NodesByID:    make(map[string]*types.Node),
		NodesByNames: make(map[string]map[string]*types.Node),
	}
}

func (r *ChainRegistry) RegisterNode(_ context.Context, node *types.Node) error {

	if !node.IsValid() {
		return errors.FromError(fmt.Errorf("invalid node")).ExtendComponent(component)
	}

	if r.NodesByNames[node.TenantID] == nil {
		r.NodesByNames[node.TenantID] = make(map[string]*types.Node)
	}

	if r.NodesByNames[node.TenantID][node.Name] != nil {
		return errors.FromError(fmt.Errorf("node tenantID=%s name=%s already exitst", node.TenantID, node.Name)).ExtendComponent(component)
	}

	node.ID = uuid.NewV4().String()
	r.NodesByNames[node.TenantID][node.Name] = node
	r.NodesByID[node.ID] = node
	return nil
}

func (r *ChainRegistry) GetNodes(_ context.Context) ([]*types.Node, error) {
	nodes := make([]*types.Node, 0)

	for _, node := range r.NodesByID {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *ChainRegistry) GetNodesByTenantID(_ context.Context, tenantID string) ([]*types.Node, error) {
	nodes := make([]*types.Node, 0)

	if tenantNodes, ok := r.NodesByNames[tenantID]; ok {
		for _, node := range tenantNodes {
			nodes = append(nodes, node)
		}
	} else {
		return nil, errors.NotFoundError("unknown tenantID=%s", tenantID).ExtendComponent(component)
	}

	return nodes, nil
}

func (r *ChainRegistry) GetNodeByName(_ context.Context, tenantID, name string) (*types.Node, error) {
	if _, ok := r.NodesByNames[tenantID]; ok {
		if node, ok := r.NodesByNames[tenantID][name]; ok {
			return node, nil
		}
		return nil, errors.NotFoundError("unknown node with tenantID=%s and name=%s", name, tenantID).ExtendComponent(component)
	}

	return nil, errors.NotFoundError("unknown node with tenantID=%s", tenantID).ExtendComponent(component)
}

func (r *ChainRegistry) GetNodeByID(_ context.Context, id string) (*types.Node, error) {
	if _, ok := r.NodesByID[id]; !ok {
		return nil, errors.FromError(fmt.Errorf("unknown node ID=%s", id)).ExtendComponent(component)
	}

	return r.NodesByID[id], nil
}

func (r *ChainRegistry) UpdateNodeByName(_ context.Context, node *types.Node) error {
	nodeToUpdate, ok := r.NodesByNames[node.TenantID][node.Name]
	if !ok {
		return errors.NotFoundError("no node found with tenantID %s and name %s", node.TenantID, node.Name).ExtendComponent(component)
	}

	nodeToUpdate.Name = node.Name
	nodeToUpdate.URLs = node.URLs
	nodeToUpdate.ListenerBackOffDuration = node.ListenerBackOffDuration
	nodeToUpdate.ListenerBlockPosition = node.ListenerBlockPosition
	nodeToUpdate.ListenerFromBlock = node.ListenerFromBlock
	nodeToUpdate.ListenerDepth = node.ListenerDepth
	currentTime := time.Now()
	nodeToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByName(_ context.Context, name, tenantID string, blockPosition int64) error {
	nodeToUpdate, ok := r.NodesByNames[tenantID][name]
	if !ok {
		return errors.NotFoundError("no node found with tenantID %s and name %s", tenantID, name).ExtendComponent(component)
	}

	nodeToUpdate.ListenerBlockPosition = blockPosition
	currentTime := time.Now()
	nodeToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateNodeByID(_ context.Context, node *types.Node) error {
	nodeToUpdate, ok := r.NodesByID[node.ID]
	if !ok {
		return errors.NotFoundError("no node found with id %s", node.ID).ExtendComponent(component)
	}

	nodeToUpdate.Name = node.Name
	nodeToUpdate.URLs = node.URLs
	nodeToUpdate.ListenerBackOffDuration = node.ListenerBackOffDuration
	nodeToUpdate.ListenerBlockPosition = node.ListenerBlockPosition
	nodeToUpdate.ListenerFromBlock = node.ListenerFromBlock
	nodeToUpdate.ListenerDepth = node.ListenerDepth
	currentTime := time.Now()
	nodeToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByID(_ context.Context, id string, blockPosition int64) error {
	nodeToUpdate, ok := r.NodesByID[id]
	if !ok {
		return errors.NotFoundError("no node found with id %s", id).ExtendComponent(component)
	}

	nodeToUpdate.ListenerBlockPosition = blockPosition
	currentTime := time.Now()
	nodeToUpdate.UpdatedAt = &currentTime

	return nil
}

func (r *ChainRegistry) DeleteNodeByName(_ context.Context, node *types.Node) error {
	if _, ok := r.NodesByNames[node.TenantID]; !ok {
		return errors.NotFoundError("no node found with tenant_id=%s", node.TenantID).ExtendComponent(component)
	}

	if _, ok := r.NodesByNames[node.TenantID][node.Name]; !ok {
		return errors.NotFoundError("no node found with tenant_id=%s and name=%s", node.TenantID, node.Name).ExtendComponent(component)
	}

	delete(r.NodesByID, r.NodesByNames[node.TenantID][node.Name].ID)
	delete(r.NodesByNames[node.TenantID], node.Name)
	return nil
}

func (r *ChainRegistry) DeleteNodeByID(_ context.Context, id string) error {
	if _, ok := r.NodesByID[id]; !ok {
		return errors.NotFoundError("no node found with id=%s", id).ExtendComponent(component)
	}

	delete(r.NodesByNames[r.NodesByID[id].TenantID], r.NodesByID[id].Name)
	delete(r.NodesByID, id)

	return nil
}