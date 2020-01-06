package types

import "context"

type ChainRegistryStore interface {
	RegisterNode(ctx context.Context, node *Node) error
	GetNodes(ctx context.Context) ([]*Node, error)
	GetNodesByTenantID(ctx context.Context, tenantID string) ([]*Node, error)
	GetNodeByName(ctx context.Context, tenantID string, name string) (*Node, error)
	GetNodeByID(ctx context.Context, ID string) (*Node, error)
	UpdateNodeByName(ctx context.Context, node *Node) error
	UpdateBlockPositionByName(ctx context.Context, name, tenantID string, blockPosition int64) error
	UpdateNodeByID(ctx context.Context, node *Node) error
	UpdateBlockPositionByID(ctx context.Context, id string, blockPosition int64) error
	DeleteNodeByName(ctx context.Context, node *Node) error
	DeleteNodeByID(ctx context.Context, ID string) error
}
