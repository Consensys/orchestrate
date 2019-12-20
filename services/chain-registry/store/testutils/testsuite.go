package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type ChainRegistryTestSuite struct {
	suite.Suite
	Store types.ChainRegistryStore
}

const (
	nodeName1 = "testNode1"
	nodeName2 = "testNode2"
	nodeName3 = "testNode3"
	tenantID1 = "tenantID1"
	tenantID2 = "tenantID2"
)

var tenantID1Nodes = map[string]*types.Node{
	nodeName1: {
		Name:                    nodeName1,
		TenantID:                tenantID1,
		URLs:                    []string{"testUrl1", "testUrl2"},
		ListenerDepth:           1,
		ListenerBlockPosition:   1,
		ListenerFromBlock:       1,
		ListenerBackOffDuration: "1s",
	},
	nodeName2: {
		Name:                    nodeName2,
		TenantID:                tenantID1,
		URLs:                    []string{"testUrl1", "testUrl2"},
		ListenerDepth:           2,
		ListenerBlockPosition:   2,
		ListenerFromBlock:       2,
		ListenerBackOffDuration: "2s",
	},
}
var tenantID2Nodes = map[string]*types.Node{
	nodeName1: {
		Name:                    nodeName1,
		TenantID:                tenantID2,
		URLs:                    []string{"testUrl1", "testUrl2"},
		ListenerDepth:           1,
		ListenerBlockPosition:   1,
		ListenerFromBlock:       1,
		ListenerBackOffDuration: "1s",
	},
	nodeName2: {
		Name:                    nodeName2,
		TenantID:                tenantID2,
		URLs:                    []string{"testUrl1", "testUrl2"},
		ListenerDepth:           2,
		ListenerBlockPosition:   2,
		ListenerFromBlock:       2,
		ListenerBackOffDuration: "2s",
	},
	nodeName3: {
		Name:                    nodeName3,
		TenantID:                tenantID2,
		URLs:                    []string{"testUrl1", "testUrl2"},
		ListenerDepth:           3,
		ListenerBlockPosition:   3,
		ListenerFromBlock:       3,
		ListenerBackOffDuration: "3s",
	},
}

var NodesSample = map[string]map[string]*types.Node{
	tenantID1: tenantID1Nodes,
	tenantID2: tenantID2Nodes,
}

func CompareNodes(t *testing.T, node1, node2 *types.Node) {
	assert.Equal(t, node1.Name, node2.Name, "Should get the same node name")
	assert.Equal(t, node1.TenantID, node2.TenantID, "Should get the same node tenantID")
	assert.Equal(t, node1.URLs, node2.URLs, "Should get the same node URLs")
	assert.Equal(t, node1.ListenerDepth, node2.ListenerDepth, "Should get the same node ListenerDepth")
	assert.Equal(t, node1.ListenerBlockPosition, node2.ListenerBlockPosition, "Should get the same node")
	assert.Equal(t, node1.ListenerFromBlock, node2.ListenerFromBlock, "Should get the same node ListenerBlockPosition")
	assert.Equal(t, node1.ListenerBackOffDuration, node2.ListenerBackOffDuration, "Should get the same node ListenerBackOffDuration")
}

func (s *ChainRegistryTestSuite) TestRegisterNode() {
	err := s.Store.RegisterNode(context.Background(), NodesSample[tenantID1][nodeName1])
	assert.NoError(s.T(), err, "Should register node properly")

	err = s.Store.RegisterNode(context.Background(), NodesSample[tenantID1][nodeName1])
	assert.Error(s.T(), err, "Should get an error violating the 'unique' constrain")
}

func (s *ChainRegistryTestSuite) TestRegisterNodes() {
	for _, nodes := range NodesSample {
		for _, node := range nodes {
			_ = s.Store.RegisterNode(context.Background(), node)
		}
	}
}

func (s *ChainRegistryTestSuite) TestGetNodes() {
	s.TestRegisterNodes()

	nodes, err := s.Store.GetNodes(context.Background())
	assert.NoError(s.T(), err, "Should get nodes without errors")
	assert.Len(s.T(), nodes, len(tenantID1Nodes)+len(tenantID2Nodes), "Should get the same number of nodes")

	for _, node := range nodes {
		CompareNodes(s.T(), node, NodesSample[node.TenantID][node.Name])
	}
}

func (s *ChainRegistryTestSuite) TestGetNodesByTenantID() {
	s.TestRegisterNodes()

	nodes, err := s.Store.GetNodesByTenantID(context.Background(), tenantID2)
	assert.NoError(s.T(), err, "Should get nodes without errors")
	assert.Len(s.T(), nodes, len(NodesSample[tenantID2]), "Should get the same number of nodes")
	for i := 0; i < len(NodesSample[tenantID2]); i++ {
		CompareNodes(s.T(), nodes[i], NodesSample[tenantID2][nodes[i].Name])
	}
}

func (s *ChainRegistryTestSuite) TestGetNodeByName() {
	s.TestRegisterNodes()

	node, err := s.Store.GetNodeByName(context.Background(), tenantID2, nodeName2)
	assert.NoError(s.T(), err, "Should get node without errors")

	CompareNodes(s.T(), node, NodesSample[tenantID2][nodeName2])
}

func (s *ChainRegistryTestSuite) TestGetNodeByID() {
	s.TestRegisterNodes()

	testNode, _ := s.Store.GetNodeByName(context.Background(), tenantID2, nodeName3)

	node, err := s.Store.GetNodeByID(context.Background(), testNode.ID)
	assert.NoError(s.T(), err, "Should get node without errors")

	CompareNodes(s.T(), node, NodesSample[tenantID2][nodeName3])
}

func (s *ChainRegistryTestSuite) TestUpdateNodeByName() {
	s.TestRegisterNodes()

	testNode := NodesSample[tenantID1][nodeName2]
	testNode.URLs = []string{"testUrl1"}
	err := s.Store.UpdateNodeByName(context.Background(), testNode)
	assert.NoError(s.T(), err, "Should get node without errors")

	node, _ := s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)
	CompareNodes(s.T(), node, testNode)
}

func (s *ChainRegistryTestSuite) TestNotFoundTenantErrorUpdateNodeByName() {
	testNode := NodesSample[tenantID1][nodeName2]
	testNode.URLs = []string{"testUrl1"}
	err := s.Store.UpdateNodeByName(context.Background(), testNode)
	assert.Error(s.T(), err, "Should get node without errors")
}

func (s *ChainRegistryTestSuite) TestNotFoundNameErrorUpdateNodeByName() {
	s.TestRegisterNodes()

	testNode := &types.Node{
		Name:     tenantID1,
		TenantID: "errorNodeName",
		URLs:     []string{"testUrl1"},
	}
	err := s.Store.UpdateNodeByName(context.Background(), testNode)
	assert.Error(s.T(), err, "Should get node without errors")
}

func (s *ChainRegistryTestSuite) TestUpdateNodeByID() {
	s.TestRegisterNodes()

	testNode, _ := s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)
	testNode.ListenerFromBlock = 10
	err := s.Store.UpdateNodeByID(context.Background(), testNode)
	assert.NoError(s.T(), err, "Should get node without errors")

	node, _ := s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)
	CompareNodes(s.T(), node, testNode)
}

func (s *ChainRegistryTestSuite) TestErrorNotFoundUpdateNodeByID() {
	s.TestRegisterNodes()

	testNode := &types.Node{
		ID:   "0d60a85e-0b90-4482-a14c-108aea2557aa",
		URLs: []string{"testUrl1"},
	}
	err := s.Store.UpdateNodeByID(context.Background(), testNode)
	assert.Error(s.T(), err, "Should update node with errors")
}

func (s *ChainRegistryTestSuite) TestDeleteNodeByName() {
	s.TestRegisterNodes()

	testNode := NodesSample[tenantID1][nodeName2]
	err := s.Store.DeleteNodeByName(context.Background(), testNode)
	assert.NoError(s.T(), err, "Should get node without errors")

	node, err := s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)
	assert.Error(s.T(), err, "Should get node without errors")
	assert.Nil(s.T(), node, "Should not get node")
}

func (s *ChainRegistryTestSuite) TestErrorNotFoundDeleteNodeByName() {
	s.TestRegisterNodes()

	testNode := &types.Node{
		Name:     tenantID1,
		TenantID: "errorNodeName",
	}
	err := s.Store.DeleteNodeByName(context.Background(), testNode)
	assert.Error(s.T(), err, "Should delete node with errors")
}

func (s *ChainRegistryTestSuite) TestDeleteNodeByID() {
	s.TestRegisterNodes()

	node, _ := s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)

	err := s.Store.DeleteNodeByID(context.Background(), node.ID)
	assert.NoError(s.T(), err, "Should get node without errors")

	node, err = s.Store.GetNodeByName(context.Background(), tenantID1, nodeName2)
	assert.Error(s.T(), err, "Should get node without errors")
	assert.Nil(s.T(), node, "Should not get node")
}

func (s *ChainRegistryTestSuite) TestErrorNotFoundDeleteNodeByID() {
	s.TestRegisterNodes()

	err := s.Store.DeleteNodeByID(context.Background(), "0d60a85e-0b90-4482-a14c-108aea2557aa")
	assert.Error(s.T(), err, "Should delete node with errors")
}
