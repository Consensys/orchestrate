package types

import "context"

type ChainRegistryStore interface {
	RegisterNode(ctx context.Context, node *Node) error
	GetNodes(ctx context.Context) ([]*Node, error)
	GetNodesByTenantID(ctx context.Context, tenantID string) ([]*Node, error)
	GetNodeByName(ctx context.Context, tenantID string, name string) (*Node, error)
	GetNodeByID(ctx context.Context, ID int) (*Node, error)
	UpdateNodeByName(ctx context.Context, node *Node) error
	UpdateNodeByID(ctx context.Context, node *Node) error
	DeleteNodeByName(ctx context.Context, node *Node) error
	DeleteNodeByID(ctx context.Context, ID int) error
}
