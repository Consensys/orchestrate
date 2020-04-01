// +build unit

package chainregistry

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var mockChain = &types.Chain{
	UUID:                    "test-chain",
	Name:                    "test",
	TenantID:                "test",
	URLs:                    []string{"test"},
	ListenerDepth:           &(&struct{ x uint64 }{0}).x,
	ListenerCurrentBlock:    &(&struct{ x uint64 }{0}).x,
	ListenerStartingBlock:   &(&struct{ x uint64 }{0}).x,
	ListenerBackOffDuration: &(&struct{ x string }{"0s"}).x,
}

type ManagerTestSuite struct {
	suite.Suite
	Manager offset.Manager
}

var mockChainRegistryClient *mocks.MockChainRegistryClient

func (s *ManagerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	mockChainRegistryClient = mocks.NewMockChainRegistryClient(ctrl)

	s.Manager = NewManager(mockChainRegistryClient)
}

func (s *ManagerTestSuite) TestManagerLastBlock() {
	updatedCurrentBlock := uint64(12)
	chain := &dynamic.Chain{
		UUID: mockChain.UUID,
	}

	mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chain.UUID).Return(mockChain, nil)
	mockChainRegistryClient.EXPECT().UpdateBlockPosition(gomock.Any(), chain.UUID, updatedCurrentBlock)

	lastBlockNumber, err := s.Manager.GetLastBlockNumber(context.Background(), chain)
	assert.Nil(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), *mockChain.ListenerCurrentBlock, lastBlockNumber, "Lastblock should be correct")

	err = s.Manager.SetLastBlockNumber(context.Background(), chain, updatedCurrentBlock)
	assert.Nil(s.T(), err, "SetLastBlockNumber should not error")
}

func (s *ManagerTestSuite) TestManagerLastIndex() {
	chain := &dynamic.Chain{
		UUID: mockChain.UUID,
	}

	lastTxIndex, err := s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")

	err = s.Manager.SetLastTxIndex(context.Background(), chain, 10, 17)
	assert.Nil(s.T(), err, "SetLastTxIndex should not error")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(17), lastTxIndex, "LastTxIndex should be correct")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 11)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")
}

func TestRegistry(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
