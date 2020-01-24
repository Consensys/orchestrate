package pg

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

// ChainRegistry is a traefik dynamic config registry based on PostgreSQL
type ChainRegistry struct {
	db *pg.DB
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(_ *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	log.Trace(q.FormattedQuery())
}

// NewChainRegistry creates a new chain registry
func NewChainRegistry(db *pg.DB) *ChainRegistry {
	db.AddQueryHook(dbLogger{})
	return &ChainRegistry{db: db}
}

// NewContractRegistryFromPGOptions creates a new pg chain registry
func NewChainRegistryPGOptions(opts *pg.Options) *ChainRegistry {
	return NewChainRegistry(pg.Connect(opts))
}

func (r *ChainRegistry) RegisterNode(ctx context.Context, node *types.Node) error {

	if node.ID == "" {
		node.ID = uuid.NewV4().String()
	}

	_, err := r.db.ModelContext(ctx, node).
		Returning("id").
		Insert()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	return nil
}

func (r *ChainRegistry) GetNodes(ctx context.Context, filters map[string]string) ([]*types.Node, error) {
	nodes := make([]*types.Node, 0)

	req := r.db.ModelContext(ctx, &nodes)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}

func (r *ChainRegistry) GetNodesByTenantID(ctx context.Context, tenantID string, filters map[string]string) ([]*types.Node, error) {
	nodes := make([]*types.Node, 0)

	req := r.db.ModelContext(ctx, &nodes).
		Where("tenant_id = ?", tenantID)
	for k, v := range filters {
		req.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := req.Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return nodes, nil
}

func (r *ChainRegistry) GetNodeByTenantIDAndNodeName(ctx context.Context, tenantID, name string) (*types.Node, error) {
	node := &types.Node{}

	err := r.db.ModelContext(ctx, node).
		Where("name = ?", name).
		Where("tenant_id = ?", tenantID).
		Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return node, nil
}

func (r *ChainRegistry) GetNodeByTenantIDAndNodeID(ctx context.Context, tenantID, id string) (*types.Node, error) {
	node := &types.Node{}

	err := r.db.ModelContext(ctx, node).
		Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return node, nil
}

func (r *ChainRegistry) GetNodeByID(ctx context.Context, id string) (*types.Node, error) {
	node := &types.Node{}

	err := r.db.ModelContext(ctx, node).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	return node, nil
}

func (r *ChainRegistry) UpdateNodeByName(ctx context.Context, node *types.Node) error {

	res, err := r.db.ModelContext(ctx, node).
		Where("tenant_id = ?", node.TenantID).
		Where("name = ?", node.Name).
		UpdateNotNull()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with tenant_id=%s and name=%s", node.TenantID, node.Name).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByName(ctx context.Context, name, tenantID string, blockPosition int64) error {
	node := &types.Node{}

	res, err := r.db.ModelContext(ctx, node).
		Set("listener_block_position = ?", blockPosition).
		Where("tenant_id = ?", tenantID).
		Where("name = ?", name).
		Update()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with tenant_id=%s and name=%s", node.TenantID, node.Name).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) UpdateNodeByID(ctx context.Context, node *types.Node) error {

	res, err := r.db.ModelContext(ctx, node).
		WherePK().
		UpdateNotNull()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with id %s", node.ID).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) UpdateBlockPositionByID(ctx context.Context, id string, blockPosition int64) error {
	node := &types.Node{}

	res, err := r.db.ModelContext(ctx, node).
		Set("listener_block_position = ?", blockPosition).
		Where("id = ?", id).
		Update()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with id %s", id).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) DeleteNodeByName(ctx context.Context, node *types.Node) error {

	res, err := r.db.ModelContext(ctx, node).
		Where("tenant_id = ?", node.TenantID).
		Where("name = ?", node.Name).
		Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with tenant_id=%s and name=%s", node.TenantID, node.Name).ExtendComponent(component)
	}
	return nil
}

func (r *ChainRegistry) DeleteNodeByID(ctx context.Context, id string) error {
	node := &types.Node{}

	res, err := r.db.ModelContext(ctx, node).
		Where("id = ?", id).
		Delete()
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	if res.RowsReturned() == 0 && res.RowsAffected() == 0 {
		return errors.NotFoundError("no node found with id %s", id).ExtendComponent(component)
	}
	return nil
}
