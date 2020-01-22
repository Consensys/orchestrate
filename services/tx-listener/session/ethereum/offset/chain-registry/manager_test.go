package chainregistry

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/testutils"
)

var MockNodesSlice = []*types.Node{
	{
		ID:                      "test-node",
		Name:                    "test",
		TenantID:                "test",
		URLs:                    []string{"test"},
		ListenerDepth:           0,
		ListenerBlockPosition:   0,
		ListenerFromBlock:       0,
		ListenerBackOffDuration: "0s",
	},
	{
		ID:                      "test-node1",
		Name:                    "test1",
		TenantID:                "test1",
		URLs:                    []string{"test1"},
		ListenerDepth:           1,
		ListenerBlockPosition:   1,
		ListenerFromBlock:       1,
		ListenerBackOffDuration: "1s",
	},
}

var MockNodesMap = map[string]*types.Node{
	"test-node":  MockNodesSlice[0],
	"test-node1": MockNodesSlice[1],
}

type Mock struct{}

func (c *Mock) GetNodeByID(_ context.Context, nodeID string) (*types.Node, error) {
	if _, ok := MockNodesMap[nodeID]; !ok {
		return nil, fmt.Errorf("test")
	}
	return MockNodesMap[nodeID], nil
}

func (c *Mock) GetNodeByTenantAndNodeName(_ context.Context, _, _ string) (*types.Node, error) {
	return nil, nil
}

func (c *Mock) GetNodeByTenantAndNodeID(_ context.Context, _, _ string) (*types.Node, error) {
	return nil, nil
}

func (c *Mock) GetNodes(_ context.Context) ([]*types.Node, error) {
	return MockNodesSlice, nil
}

func (c *Mock) UpdateBlockPosition(_ context.Context, nodeID string, blockNumber int64) error {
	if _, ok := MockNodesMap[nodeID]; !ok {
		return fmt.Errorf("test")
	}
	MockNodesMap[nodeID].ListenerBlockPosition = blockNumber
	return nil
}

type ManagerTestSuite struct {
	testutils.OffsetManagerTestSuite
}

func (s *ManagerTestSuite) SetupTest() {
	s.Manager = NewManager(&Mock{})
}

func TestRegistry(t *testing.T) {
	s := new(ManagerTestSuite)
	suite.Run(t, s)
}
